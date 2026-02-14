# AWS Lambda Deployment - Command Checklist

Complete list of commands needed to deploy the dev environment to AWS using Terraform.

---

## Prerequisites Verification

```bash
# 1. Check AWS CLI is installed and configured
aws --version
aws sts get-caller-identity

# 2. Check Terraform is installed
terraform --version

# 3. Check Go is installed
go version
```

---

## Step 1: Create S3 Backend for Terraform State

```bash
# Create S3 bucket
aws s3 mb s3://dora-metrics-terraform-state-dev --region ap-southeast-1

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket dora-metrics-terraform-state-dev \
  --versioning-configuration Status=Enabled \
  --region ap-southeast-1

# Enable encryption
aws s3api put-bucket-encryption \
  --bucket dora-metrics-terraform-state-dev \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }' \
  --region ap-southeast-1

# Block public access
aws s3api put-public-access-block \
  --bucket dora-metrics-terraform-state-dev \
  --public-access-block-configuration \
    "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true" \
  --region ap-southeast-1
```

---

## Step 2: Create DynamoDB Table for State Locking

```bash
aws dynamodb create-table \
  --table-name dora-metrics-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region ap-southeast-1
```

---

## Step 3: Store GitHub PAT in Secrets Manager

```bash
# Replace YOUR_GITHUB_PAT with your actual token
aws secretsmanager create-secret \
  --name dora-metrics/dev/github-pat \
  --description "GitHub Personal Access Token for DORA metrics collector" \
  --secret-string "YOUR_GITHUB_PAT" \
  --region ap-southeast-1
```

---

## Step 4: Build Lambda Deployment Package

```bash
# From project root directory
cd /Users/lamdo/Projects/go-github-tracker

# Make build script executable
chmod +x scripts/build-lambda.sh

# Build ARM64 binary and create deployment package
./scripts/build-lambda.sh
```

**Expected output**: `lambda-deployment.zip` created in project root

---

## Step 5: Configure Terraform Variables

```bash
# Navigate to terraform dev directory
cd terraform/dev

# Copy example tfvars
cp terraform.tfvars.example terraform.tfvars

# Edit terraform.tfvars with your values
# Use your preferred editor (vim, nano, code, etc.)
vim terraform.tfvars
```

**Required edits in `terraform.tfvars`**:
```hcl
github_repositories = "your-org/your-repo1,your-org/your-repo2"

team_config_json = <<EOF
[
  {
    "name": "Your Team Name",
    "members": [
      {"username": "github-user1", "allocation": 1.0},
      {"username": "github-user2", "allocation": 1.0}
    ]
  }
]
EOF
```

---

## Step 6: Initialize Terraform

```bash
# Still in terraform/dev directory
terraform init
```

**Expected output**:
- Downloads AWS provider
- Configures S3 backend
- Initializes modules

---

## Step 7: Validate Terraform Configuration

```bash
# Validate syntax
terraform validate

# Format files
terraform fmt -recursive
```

---

## Step 8: Review Terraform Plan

```bash
# Generate and review execution plan
terraform plan

# Optional: Save plan to file for review
terraform plan -out=tfplan
```

**Review the plan carefully**:
- VPC and subnets
- RDS PostgreSQL instance
- Lambda function
- EventBridge rule
- Secrets Manager secrets
- CloudWatch alarms
- IAM roles and policies

---

## Step 9: Apply Terraform Configuration

```bash
# Apply the configuration
terraform apply

# Or if you saved the plan:
terraform apply tfplan
```

**This will create**:
- VPC with public/private subnets
- NAT Gateway
- RDS PostgreSQL (db.t4g.micro)
- Lambda function
- EventBridge schedule
- Secrets Manager secrets
- CloudWatch log groups and alarms

**Expected duration**: 10-15 minutes (RDS takes the longest)

---

## Step 10: View Terraform Outputs

```bash
# View all outputs
terraform output

# View specific output
terraform output lambda_function_name
terraform output rds_endpoint
```

---

## Step 11: Update Lambda Function Code

```bash
# Go back to project root
cd ../..

# Get Lambda function name from Terraform
FUNCTION_NAME=$(cd terraform/dev && terraform output -raw lambda_function_name)

# Update Lambda function code
aws lambda update-function-code \
  --function-name $FUNCTION_NAME \
  --zip-file fileb://lambda-deployment.zip \
  --region ap-southeast-1
```

---

## Step 12: Test Lambda Function

```bash
# Invoke Lambda function manually
aws lambda invoke \
  --function-name $FUNCTION_NAME \
  --payload '{}' \
  --region ap-southeast-1 \
  response.json

# View response
cat response.json
```

---

## Step 13: Check CloudWatch Logs

```bash
# Get log group name
LOG_GROUP=$(cd terraform/dev && terraform output -raw cloudwatch_log_group)

# Tail logs in real-time
aws logs tail $LOG_GROUP --follow --region ap-southeast-1

# Or view recent logs
aws logs tail $LOG_GROUP --since 1h --region ap-southeast-1
```

---

## Step 14: Verify Database Connection

```bash
# Get RDS endpoint
RDS_ENDPOINT=$(cd terraform/dev && terraform output -raw rds_endpoint)

# Get database credentials from Secrets Manager
DB_SECRET_ARN=$(cd terraform/dev && terraform output -raw db_credentials_secret_arn)

aws secretsmanager get-secret-value \
  --secret-id $DB_SECRET_ARN \
  --region ap-southeast-1 \
  --query SecretString \
  --output text

# Connect to database (you'll need the password from above)
psql -h $RDS_ENDPOINT -U postgres -d dora_metrics
```

**In psql**:
```sql
-- Check tables
\dt

-- Check data
SELECT COUNT(*) FROM pr_metrics;
SELECT * FROM view_team_velocity LIMIT 10;

-- Exit
\q
```

---

## Step 15: Verify EventBridge Schedule

```bash
# Check EventBridge rule
aws events describe-rule \
  --name dora-metrics-dev-schedule \
  --region ap-southeast-1

# List targets
aws events list-targets-by-rule \
  --rule dora-metrics-dev-schedule \
  --region ap-southeast-1
```

---

## Step 16: Monitor CloudWatch Alarms

```bash
# List all alarms
aws cloudwatch describe-alarms \
  --alarm-name-prefix dora-metrics-dev \
  --region ap-southeast-1

# Check alarm state
aws cloudwatch describe-alarms \
  --alarm-names \
    dora-metrics-dev-lambda-errors \
    dora-metrics-dev-lambda-duration \
    dora-metrics-dev-lambda-throttles \
  --region ap-southeast-1
```

---

## Verification Checklist

After deployment, verify:

- [ ] S3 bucket created for Terraform state
- [ ] DynamoDB table created for state locking
- [ ] GitHub PAT stored in Secrets Manager
- [ ] Lambda deployment package built
- [ ] Terraform initialized successfully
- [ ] Terraform plan reviewed
- [ ] Terraform apply completed
- [ ] Lambda function deployed
- [ ] Lambda function invoked successfully
- [ ] CloudWatch logs showing execution
- [ ] RDS database accessible
- [ ] EventBridge rule created
- [ ] CloudWatch alarms configured

---

## Cost Monitoring

```bash
# Check current month costs (requires Cost Explorer API)
aws ce get-cost-and-usage \
  --time-period Start=$(date -u +%Y-%m-01),End=$(date -u +%Y-%m-%d) \
  --granularity MONTHLY \
  --metrics BlendedCost \
  --filter file://cost-filter.json \
  --region ap-southeast-1
```

---

## Cleanup (When Needed)

```bash
# Destroy all resources
cd terraform/dev
terraform destroy

# Delete S3 bucket (after destroying resources)
aws s3 rb s3://dora-metrics-terraform-state-dev --force --region ap-southeast-1

# Delete DynamoDB table
aws dynamodb delete-table \
  --table-name dora-metrics-terraform-locks \
  --region ap-southeast-1

# Delete secrets
aws secretsmanager delete-secret \
  --secret-id dora-metrics/dev/github-pat \
  --force-delete-without-recovery \
  --region ap-southeast-1
```

---

## Summary

**Total commands**: ~30 commands across 16 steps

**Estimated time**:
- Setup (Steps 1-5): 10 minutes
- Terraform deployment (Steps 6-9): 15 minutes
- Verification (Steps 10-16): 10 minutes
- **Total**: ~35 minutes

**Estimated cost**: ~$17/month for dev environment
