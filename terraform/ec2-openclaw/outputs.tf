output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.openclaw.id
}

output "public_ip" {
  description = "The public IP of the EC2 instance"
  value       = aws_instance.openclaw.public_ip
}

output "private_ip" {
  description = "The private IP of the EC2 instance"
  value       = aws_instance.openclaw.private_ip
}

output "security_group_id" {
  description = "The ID of the security group"
  value       = aws_security_group.openclaw.id
}
