variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "memory_size" {
  description = "Lambda memory size in MB"
  type        = number
}

variable "timeout" {
  description = "Lambda timeout in seconds"
  type        = number
}

variable "deployment_package_path" {
  description = "Path to Lambda deployment package"
  type        = string
  default     = "../../lambda-deployment.zip"
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for Lambda"
  type        = list(string)
}

variable "security_group_ids" {
  description = "Security group IDs for Lambda"
  type        = list(string)
}

variable "db_host" {
  description = "Database host"
  type        = string
}

variable "db_name" {
  description = "Database name"
  type        = string
}

variable "db_secret_arn" {
  description = "ARN of database credentials secret"
  type        = string
}

variable "github_pat_secret_arn" {
  description = "ARN of GitHub PAT secret"
  type        = string
}

variable "github_repositories" {
  description = "Comma-separated list of repositories"
  type        = string
}

variable "team_config_json" {
  description = "Team configuration JSON"
  type        = string
}
