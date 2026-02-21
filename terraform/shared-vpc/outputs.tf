output "vpc_id" {
  description = "Shared VPC ID — paste into dev/terraform.tfvars as vpc_id"
  value       = module.vpc.vpc_id
}

output "private_subnet_ids" {
  description = "Private subnet IDs (for Lambda and RDS) — paste into dev/terraform.tfvars"
  value       = module.vpc.private_subnet_ids
}

output "public_subnet_ids" {
  description = "Public subnet IDs — paste into dev/terraform.tfvars if needed"
  value       = module.vpc.public_subnet_ids
}

output "nat_gateway_id" {
  description = "NAT Gateway ID"
  value       = module.vpc.nat_gateway_id
}
