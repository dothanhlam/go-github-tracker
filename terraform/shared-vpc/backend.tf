terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "dora-metrics-terraform-state-dev"
    key            = "shared-vpc/terraform.tfstate"
    region         = "ap-southeast-1"
    dynamodb_table = "dora-metrics-terraform-locks"
    encrypt        = true
  }
}

provider "aws" {
  region = var.aws_region
}
