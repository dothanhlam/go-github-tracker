# Data source for latest Amazon Linux 2023 AMI
data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# Security Group for OpenClaw
resource "aws_security_group" "openclaw" {
  name        = "${var.project_name}-${var.environment}-openclaw-sg"
  description = "Security group for OpenClaw bot"
  vpc_id      = var.vpc_id

  # OpenClaw Control UI / Gateway
  ingress {
    from_port   = 18789
    to_port     = 18789
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Consider restricting this in production
    description = "OpenClaw Gateway"
  }

  # SSH Access
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.allowed_ssh_cidr
    description = "SSH Access"
  }

  # Outbound everything
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-openclaw-sg"
  }
}

# IAM Role for SSM Access
resource "aws_iam_role" "openclaw" {
  name = "${var.project_name}-${var.environment}-openclaw-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })

  tags = {
    Name = "${var.project_name}-${var.environment}-openclaw-role"
  }
}

resource "aws_iam_role_policy_attachment" "ssm_core" {
  role       = aws_iam_role.openclaw.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "openclaw" {
  name = "${var.project_name}-${var.environment}-openclaw-profile"
  role = aws_iam_role.openclaw.name
}

# EC2 Instance
resource "aws_instance" "openclaw" {
  ami                    = var.ami_id != null ? var.ami_id : data.aws_ami.amazon_linux_2023.id
  instance_type               = var.instance_type
  subnet_id                   = var.subnet_id
  vpc_security_group_ids      = [aws_security_group.openclaw.id]
  iam_instance_profile        = aws_iam_instance_profile.openclaw.name
  associate_public_ip_address = var.associate_public_ip_address

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "optional" # Allow v1 and v2 for maximum compatibility
    http_put_response_hop_limit = 2
  }

  user_data = <<-EOF
#!/bin/bash
# Enable logging
exec > >(tee /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1

echo "Starting user_data script..."

# Ensure SSM Agent is installed and running (unconditionally)
echo "Ensuring SSM Agent is installed..."
dnf install -y amazon-ssm-agent
systemctl enable amazon-ssm-agent
systemctl stop amazon-ssm-agent
systemctl start amazon-ssm-agent
echo "SSM Agent status:"
systemctl status amazon-ssm-agent --no-pager

# Install Node.js 22
echo "Installing Node.js 22..."
curl -fsSL https://rpm.nodesource.com/setup_22.x | bash -
dnf install -y nodejs

# Install pnpm
echo "Installing pnpm..."
npm install -g pnpm

echo "user_data script finished."
EOF

  tags = {
    Name = "${var.project_name}-${var.environment}-openclaw"
  }

  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }
}

# EC2 Auto Stop/Start Scheduler (10 PM stop, 10 AM start — ICT timezone)
module "ec2_scheduler" {
  source = "../modules/ec2-scheduler"

  project_name = var.project_name
  environment  = var.environment
  instance_id  = aws_instance.openclaw.id
  instance_arn = aws_instance.openclaw.arn
}
