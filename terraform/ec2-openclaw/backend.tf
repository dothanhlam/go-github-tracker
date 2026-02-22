terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # Local backend for now, or you can configure S3 if needed
  # backend "s3" { ... }
}

provider "aws" {
  region = var.aws_region
}
