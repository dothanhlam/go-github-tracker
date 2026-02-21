# Secrets Manager
module "secrets" {
  source = "../modules/secrets"

  project_name = var.project_name
  environment  = var.environment
}

# RDS PostgreSQL
module "rds" {
  source = "../modules/rds"

  project_name      = var.project_name
  environment       = var.environment
  db_name           = var.db_name
  instance_class    = var.db_instance_class
  allocated_storage = var.db_allocated_storage
  vpc_id            = var.vpc_id
  subnet_ids        = var.private_subnet_ids
  db_secret_arn     = module.secrets.db_credentials_secret_arn
  db_password       = module.secrets.db_password
}

# Lambda Function
module "lambda" {
  source = "../modules/lambda"

  project_name          = var.project_name
  environment           = var.environment
  memory_size           = var.lambda_memory
  timeout               = var.lambda_timeout
  vpc_id                = var.vpc_id
  subnet_ids            = var.private_subnet_ids
  security_group_ids    = [module.rds.lambda_security_group_id]
  db_host               = module.rds.db_endpoint
  db_name               = var.db_name
  db_secret_arn         = module.secrets.db_credentials_secret_arn
  github_pat_secret_arn = module.secrets.github_pat_secret_arn
  github_repositories   = var.github_repositories
  team_config_json      = var.team_config_json
}

# EventBridge Scheduler
module "eventbridge" {
  source = "../modules/eventbridge"

  project_name         = var.project_name
  environment          = var.environment
  schedule_expression  = var.schedule_expression
  lambda_function_arn  = module.lambda.function_arn
  lambda_function_name = module.lambda.function_name
}

# CloudWatch Monitoring
module "monitoring" {
  source = "../modules/monitoring"

  project_name         = var.project_name
  environment          = var.environment
  lambda_function_name = module.lambda.function_name
}

# RDS Auto Stop/Start Scheduler (10 PM stop, 10 AM start — ICT timezone)
module "rds_scheduler" {
  source = "../modules/rds-scheduler"

  project_name           = var.project_name
  environment            = var.environment
  db_instance_identifier = module.rds.db_instance_identifier
  db_instance_arn        = module.rds.db_instance_arn
}

# GitHub Actions OIDC — allows GitHub Actions to deploy Lambda without stored AWS keys
module "github_oidc" {
  source = "../modules/github-oidc"

  project_name = var.project_name
  github_org   = var.github_org
  github_repo  = var.github_repo
}
