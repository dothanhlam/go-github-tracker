output "db_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.main.endpoint
}

output "db_address" {
  description = "RDS instance address"
  value       = aws_db_instance.main.address
}

output "db_port" {
  description = "RDS instance port"
  value       = aws_db_instance.main.port
}

output "lambda_security_group_id" {
  description = "Security group ID for Lambda to access RDS"
  value       = aws_security_group.lambda.id
}

output "db_instance_identifier" {
  description = "RDS DB instance identifier"
  value       = aws_db_instance.main.identifier
}

output "db_instance_arn" {
  description = "RDS DB instance ARN"
  value       = aws_db_instance.main.arn
}
