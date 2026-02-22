output "scheduler_group_name" {
  description = "Name of the EventBridge Scheduler group"
  value       = aws_scheduler_schedule_group.ec2.name
}
