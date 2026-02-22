variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "shared"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "openclaw"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.small"
}

variable "ami_id" {
  description = "AMI ID for the EC2 instance (Amazon Linux 2023 recommended)"
  type        = string
  default     = null
}

variable "vpc_id" {
  description = "VPC ID where the host will be created"
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID where the host will be created"
  type        = string
}

variable "allowed_ssh_cidr" {
  description = "Allowed CIDR blocks for SSH access"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "associate_public_ip_address" {
  description = "Whether to associate a public IP address with the instance"
  type        = bool
  default     = false
}
