#!/bin/bash
set -e

echo "⚠️ This script will destroy the 'dev' Terraform environment and permanently delete the 'rds-scheduler' module."
echo "Press Enter to continue or Ctrl+C to abort..."
read -r

# 1. Navigate to the dev environment and destroy the infrastructure
echo "=================================================="
echo "Step 1: Destroying the Dev environment"
echo "=================================================="
cd /Users/lamdo/Projects/go-github-tracker/terraform/dev
terraform destroy

# 2. Once destroyed, remove the rds_scheduler module from main.tf
echo "=================================================="
echo "Step 2: Removing rds_scheduler configuration from main.tf"
echo "=================================================="
# Using sed to delete the module block and its preceding comment
sed -i.bak '/# RDS Auto Stop\/Start Scheduler/,/^}/d' main.tf
rm -f main.tf.bak

# 3. Remove the modules/rds-scheduler directory
echo "=================================================="
echo "Step 3: Deleting the rds-scheduler module folder"
echo "=================================================="
cd /Users/lamdo/Projects/go-github-tracker/terraform
rm -rf modules/rds-scheduler

echo "=================================================="
echo "✅ Cleanup complete!"
echo "The dev environment has been destroyed and the rds-scheduler module is removed."
