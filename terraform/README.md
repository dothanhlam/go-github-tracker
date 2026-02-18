# Terraform AWS Deployment

This directory contains Terraform infrastructure code for deploying the DORA metrics collector to AWS Lambda.

## Architecture

- **Lambda Function**: ARM64 Go runtime, runs every 4 hours
- **RDS PostgreSQL**: Database for metrics storage
- **EventBridge**: Scheduled trigger
- **Secrets Manager**: GitHub PAT and DB credentials
- **CloudWatch**: Logging and monitoring

## Directory Structure

```
terraform/
├── dev/                    # Dev environment
│   ├── main.tf            # Main configuration
│   ├── variables.tf       # Input variables
│   ├── outputs.tf         # Output values
│   ├── backend.tf         # S3 backend config
│   └── terraform.tfvars   # Variable values (gitignored)
└── modules/               # Reusable modules
    ├── vpc/              # VPC, subnets, NAT gateway
    ├── secrets/          # Secrets Manager
    ├── rds/              # PostgreSQL database
    ├── lambda/           # Lambda function
    ├── eventbridge/      # Scheduled triggers
    └── monitoring/       # CloudWatch alarms
```

## Prerequisites

1. **AWS CLI** configured with credentials
2. **Terraform** >= 1.0 installed
3. **Go** 1.21+ for building Lambda binary

## Initial Setup

### 1. Create S3 Backend

```bash
# Create S3 bucket for Terraform state
aws s3 mb s3://dora-metrics-terraform-state-dev --region ap-southeast-1

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket dora-metrics-terraform-state-dev \
  --versioning-configuration Status=Enabled

# Enable encryption
aws s3api put-bucket-encryption \
  --bucket dora-metrics-terraform-state-dev \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'
```

### 2. Create DynamoDB Table for State Locking

```bash
aws dynamodb create-table \
  --table-name dora-metrics-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region ap-southeast-1
```

### 3. Build Lambda Deployment Package

```bash
# From project root
./scripts/build-lambda.sh
```

This creates `lambda-deployment.zip` with the ARM64 binary.

### 4. Configure Variables

```bash
cd terraform/dev
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
```

**Required variables**:
- `github_repositories` - Repositories to track
- `team_config_json` - Team configuration

**Optional variables**:
- `COLLECTION_LOOKBACK_DAYS` - Days to look back for PR collection (default: 90)

### 5. Set GitHub PAT (One-time)

```bash
# Store GitHub PAT in Secrets Manager
aws secretsmanager create-secret \
  --name dora-metrics/dev/github-pat \
  --secret-string "your-github-pat-here" \
  --region ap-southeast-1
```

## Deployment

### Deploy Infrastructure

```bash
cd terraform/dev

# Initialize Terraform
terraform init

# Review plan
terraform plan

# Apply changes
terraform apply
```

### Update Lambda Code

After infrastructure is deployed, update Lambda function code:

```bash
# From project root
./scripts/build-lambda.sh

# Update Lambda function
aws lambda update-function-code \
  --function-name dora-metrics-dev-collector \
  --zip-file fileb://lambda-deployment.zip \
  --region ap-southeast-1
```

## Verification

### Check Terraform Outputs

```bash
cd terraform/dev
terraform output
```

### Test Lambda Function

```bash
aws lambda invoke \
  --function-name dora-metrics-dev-collector \
  --payload '{}' \
  --region ap-southeast-1 \
  response.json

cat response.json
```

### View Logs

```bash
# Tail logs
aws logs tail /aws/lambda/dora-metrics-dev-collector --follow

# Or use AWS Console
# CloudWatch > Log groups > /aws/lambda/dora-metrics-dev-collector
```

### Check Database

```bash
# Get RDS endpoint
terraform output rds_endpoint

# Connect to database
psql -h <rds-endpoint> -U postgres -d dora_metrics

# Check data
SELECT COUNT(*) FROM pr_metrics;
SELECT * FROM view_team_velocity LIMIT 10;
```

### Verify Schedule

```bash
# Check EventBridge rule
aws events describe-rule \
  --name dora-metrics-dev-schedule \
  --region ap-southeast-1

# List recent invocations
aws lambda get-function \
  --function-name dora-metrics-dev-collector \
  --region ap-southeast-1
```

## Cost Estimation

**Dev environment** (~$17/month):
- Lambda: ~$0.50
- RDS (db.t4g.micro): ~$15
- Secrets Manager: ~$0.80
- CloudWatch: ~$0.50
- S3/DynamoDB: ~$0.20

## Troubleshooting

### Lambda can't connect to RDS

- Check security groups allow Lambda → RDS on port 5432
- Verify Lambda is in private subnets
- Check NAT gateway is working

### Secrets not accessible

- Verify IAM role has `secretsmanager:GetSecretValue` permission
- Check secret ARNs in Lambda environment variables

### Timeout errors

- Increase Lambda timeout (current: 300s)
- Reduce `COLLECTION_LOOKBACK_DAYS` to limit PR collection (e.g., 30 or 60 days)
- Check database query performance
- Review CloudWatch logs for bottlenecks

## Cleanup

To destroy all resources:

```bash
cd terraform/dev
terraform destroy
```

**Warning**: This will delete the RDS database and all data!

## Next Steps

1. Monitor for 1 week
2. Review CloudWatch metrics
3. Optimize Lambda memory/timeout if needed
4. Plan production environment promotion
