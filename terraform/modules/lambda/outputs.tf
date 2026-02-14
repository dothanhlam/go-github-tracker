output "function_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.collector.function_name
}

output "function_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.collector.arn
}

output "function_invoke_arn" {
  description = "Lambda function invoke ARN"
  value       = aws_lambda_function.collector.invoke_arn
}

output "log_group_name" {
  description = "CloudWatch log group name"
  value       = aws_cloudwatch_log_group.lambda.name
}
