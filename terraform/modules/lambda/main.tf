# IAM Role for Lambda
resource "aws_iam_role" "lambda" {
  name = "${var.project_name}-${var.environment}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  tags = {
    Name = "${var.project_name}-${var.environment}-lambda-role"
  }
}

# Attach AWS managed policy for VPC execution
resource "aws_iam_role_policy_attachment" "lambda_vpc" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# Custom policy for Secrets Manager access
resource "aws_iam_policy" "lambda_secrets" {
  name        = "${var.project_name}-${var.environment}-lambda-secrets"
  description = "Allow Lambda to read secrets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = [
          var.db_secret_arn,
          var.github_pat_secret_arn
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_secrets" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda_secrets.arn
}

# Lambda Function
resource "aws_lambda_function" "collector" {
  filename      = var.deployment_package_path
  function_name = "${var.project_name}-${var.environment}-collector"
  role          = aws_iam_role.lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2"
  architectures = ["arm64"]

  memory_size = var.memory_size
  timeout     = var.timeout

  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = var.security_group_ids
  }

  environment {
    variables = {
      DB_DRIVER           = "postgres"
      DB_HOST             = var.db_host
      DB_NAME             = var.db_name
      DB_SECRET_ARN       = var.db_secret_arn
      GITHUB_PAT_SECRET_ARN = var.github_pat_secret_arn
      REPOSITORIES        = var.github_repositories
      TEAM_CONFIG_JSON    = var.team_config_json
    }
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-collector"
  }

  lifecycle {
    ignore_changes = [filename, source_code_hash]
  }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${aws_lambda_function.collector.function_name}"
  retention_in_days = 14

  tags = {
    Name = "${var.project_name}-${var.environment}-lambda-logs"
  }
}
