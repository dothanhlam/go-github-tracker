output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = module.lambda.function_name
}

output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = module.lambda.function_arn
}

output "rds_endpoint" {
  description = "RDS instance endpoint"
  value       = module.rds.db_endpoint
}

output "rds_database_name" {
  description = "RDS database name"
  value       = var.db_name
}

output "eventbridge_rule_name" {
  description = "EventBridge rule name"
  value       = module.eventbridge.rule_name
}

output "github_pat_secret_arn" {
  description = "ARN of GitHub PAT secret"
  value       = module.secrets.github_pat_secret_arn
}

output "db_credentials_secret_arn" {
  description = "ARN of database credentials secret"
  value       = module.secrets.db_credentials_secret_arn
}

output "cloudwatch_log_group" {
  description = "CloudWatch log group for Lambda"
  value       = module.monitoring.log_group_name
}
