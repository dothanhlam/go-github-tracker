# IAM Role for EventBridge Scheduler to call RDS APIs
resource "aws_iam_role" "rds_scheduler" {
  name = "${var.project_name}-${var.environment}-rds-scheduler-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "scheduler.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })

  tags = {
    Name = "${var.project_name}-${var.environment}-rds-scheduler-role"
  }
}

resource "aws_iam_role_policy" "rds_scheduler" {
  name = "${var.project_name}-${var.environment}-rds-scheduler-policy"
  role = aws_iam_role.rds_scheduler.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "rds:StopDBInstance",
          "rds:StartDBInstance"
        ]
        Resource = var.db_instance_arn
      }
    ]
  })
}

# Scheduler Group
resource "aws_scheduler_schedule_group" "rds" {
  name = "${var.project_name}-${var.environment}-rds-scheduler"

  tags = {
    Name = "${var.project_name}-${var.environment}-rds-scheduler"
  }
}

# Stop RDS at 10:00 PM (Asia/Ho_Chi_Minh)
resource "aws_scheduler_schedule" "rds_stop" {
  name       = "${var.project_name}-${var.environment}-rds-stop"
  group_name = aws_scheduler_schedule_group.rds.name

  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression          = var.stop_schedule
  schedule_expression_timezone = var.timezone

  target {
    arn      = "arn:aws:scheduler:::aws-sdk:rds:stopDBInstance"
    role_arn = aws_iam_role.rds_scheduler.arn

    input = jsonencode({
      DbInstanceIdentifier = var.db_instance_identifier
    })

    retry_policy {
      maximum_retry_attempts = 3
    }
  }
}

# Start RDS at 10:00 AM (Asia/Ho_Chi_Minh)
resource "aws_scheduler_schedule" "rds_start" {
  name       = "${var.project_name}-${var.environment}-rds-start"
  group_name = aws_scheduler_schedule_group.rds.name

  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression          = var.start_schedule
  schedule_expression_timezone = var.timezone

  target {
    arn      = "arn:aws:scheduler:::aws-sdk:rds:startDBInstance"
    role_arn = aws_iam_role.rds_scheduler.arn

    input = jsonencode({
      DbInstanceIdentifier = var.db_instance_identifier
    })

    retry_policy {
      maximum_retry_attempts = 3
    }
  }
}
