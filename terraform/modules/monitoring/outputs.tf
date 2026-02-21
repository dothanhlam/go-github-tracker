output "log_group_name" {
  description = "CloudWatch log group name"
  value       = "/aws/lambda/${var.lambda_function_name}"
}

output "error_alarm_arn" {
  description = "Error alarm ARN"
  value       = aws_cloudwatch_metric_alarm.lambda_errors.arn
}

output "duration_alarm_arn" {
  description = "Duration alarm ARN"
  value       = aws_cloudwatch_metric_alarm.lambda_duration.arn
}

output "throttle_alarm_arn" {
  description = "Throttle alarm ARN"
  value       = aws_cloudwatch_metric_alarm.lambda_throttles.arn
}
