variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "ap-southeast-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "dora-metrics"
}

variable "lambda_memory" {
  description = "Lambda function memory in MB"
  type        = number
  default     = 512
}

variable "lambda_timeout" {
  description = "Lambda function timeout in seconds"
  type        = number
  default     = 300
}

variable "schedule_expression" {
  description = "EventBridge schedule expression"
  type        = string
  default     = "rate(4 hours)"
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t4g.micro"
}

variable "db_allocated_storage" {
  description = "RDS allocated storage in GB"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "dora_metrics"
}

variable "github_repositories" {
  description = "Comma-separated list of GitHub repositories to track"
  type        = string
}

variable "team_config_json" {
  description = "JSON configuration for teams"
  type        = string
}

variable "vpc_id" {
  description = "Shared VPC ID (from terraform/shared-vpc output)"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs from the shared VPC (for Lambda and RDS)"
  type        = list(string)
}

variable "github_org" {
  description = "GitHub organisation or username (e.g. 'dothanhlam')"
  type        = string
  default     = "dothanhlam"
}

variable "github_repo" {
  description = "GitHub repository name (e.g. 'go-github-tracker')"
  type        = string
  default     = "go-github-tracker"
}
