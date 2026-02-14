#!/bin/bash
set -e

echo "🔍 Verifying Terraform Configuration..."
echo ""

cd terraform/dev

echo "📦 Step 1: Initializing Terraform..."
terraform init
echo "✅ Initialization complete"
echo ""

echo "🔍 Step 2: Validating syntax..."
terraform validate
echo "✅ Validation passed"
echo ""

echo "📝 Step 3: Checking format..."
if terraform fmt -check -recursive; then
    echo "✅ Format check passed"
else
    echo "⚠️  Files need formatting. Run: terraform fmt -recursive"
fi
echo ""

echo "📋 Step 4: Generating plan..."
terraform plan
echo ""

echo "✨ Verification complete!"
echo ""
echo "Next steps:"
echo "1. Review the plan output above"
echo "2. Verify region is ap-southeast-1"
echo "3. Check resource counts match expectations"
echo "4. When ready: terraform apply"