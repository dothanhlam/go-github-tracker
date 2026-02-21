# Shared VPC - created once, used by all environments and projects
module "vpc" {
  source = "../modules/vpc"

  vpc_name = var.vpc_name
  vpc_cidr = var.vpc_cidr
}
