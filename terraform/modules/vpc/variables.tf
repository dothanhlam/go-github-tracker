variable "vpc_name" {
  description = "Name prefix for the VPC and all networking resources (e.g. 'personal-ap-southeast-1')"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}
