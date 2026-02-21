output "deploy_role_arn" {
  description = "ARN of the IAM role for GitHub Actions to assume — set as AWS_DEPLOY_ROLE_ARN secret in GitHub"
  value       = aws_iam_role.github_actions_deploy.arn
}
