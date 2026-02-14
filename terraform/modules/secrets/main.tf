# GitHub PAT Secret
resource "aws_secretsmanager_secret" "github_pat" {
  name        = "${var.project_name}/${var.environment}/github-pat"
  description = "GitHub Personal Access Token for ${var.project_name}"

  tags = {
    Name = "${var.project_name}-${var.environment}-github-pat"
  }
}

# Placeholder for GitHub PAT (must be set manually or via CI/CD)
resource "aws_secretsmanager_secret_version" "github_pat" {
  secret_id     = aws_secretsmanager_secret.github_pat.id
  secret_string = var.github_pat_value != "" ? var.github_pat_value : "PLACEHOLDER_SET_VIA_AWS_CONSOLE"

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# Database Credentials Secret
resource "aws_secretsmanager_secret" "db_credentials" {
  name        = "${var.project_name}/${var.environment}/db-credentials"
  description = "Database credentials for ${var.project_name}"

  tags = {
    Name = "${var.project_name}-${var.environment}-db-credentials"
  }
}

# Generate random password for database
resource "random_password" "db_password" {
  length  = 32
  special = true
}

# Store database credentials as JSON
resource "aws_secretsmanager_secret_version" "db_credentials" {
  secret_id = aws_secretsmanager_secret.db_credentials.id
  secret_string = jsonencode({
    username = "postgres"
    password = random_password.db_password.result
  })
}
