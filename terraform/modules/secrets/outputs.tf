output "github_pat_secret_arn" {
  description = "ARN of GitHub PAT secret"
  value       = data.aws_secretsmanager_secret.github_pat.arn
}

output "db_credentials_secret_arn" {
  description = "ARN of database credentials secret"
  value       = aws_secretsmanager_secret.db_credentials.arn
}

output "db_password" {
  description = "Generated database password"
  value       = random_password.db_password.result
  sensitive   = true
}
