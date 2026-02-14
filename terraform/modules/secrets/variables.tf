terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "github_pat_value" {
  description = "GitHub PAT value (optional, can be set manually)"
  type        = string
  default     = ""
  sensitive   = true
}
