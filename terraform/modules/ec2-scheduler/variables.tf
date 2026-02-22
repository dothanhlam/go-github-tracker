variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "instance_id" {
  description = "EC2 instance ID to stop/start"
  type        = string
}

variable "instance_arn" {
  description = "EC2 instance ARN (used to scope IAM permissions)"
  type        = string
}

variable "stop_schedule" {
  description = "Cron expression (in local timezone) to stop the EC2 instance"
  type        = string
  default     = "cron(0 22 * * ? *)" # 10:00 PM
}

variable "start_schedule" {
  description = "Cron expression (in local timezone) to start the EC2 instance"
  type        = string
  default     = "cron(0 10 * * ? *)" # 10:00 AM
}

variable "timezone" {
  description = "Timezone for cron expressions"
  type        = string
  default     = "Asia/Ho_Chi_Minh"
}
