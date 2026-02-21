output "stop_schedule_arn" {
  description = "ARN of the RDS stop schedule"
  value       = aws_scheduler_schedule.rds_stop.arn
}

output "start_schedule_arn" {
  description = "ARN of the RDS start schedule"
  value       = aws_scheduler_schedule.rds_start.arn
}
