#!/bin/bash
set -e

# Configuration
ROLE_NAME="myna-test-lambda-role"
LAMBDA_DIR="./lambda"
RUNTIME="nodejs20.x"

# ANSI Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting deployment of Myna test lambdas...${NC}"

# 1. Create IAM Role if not exists
echo "Checking IAM Role: $ROLE_NAME"
if aws iam get-role --role-name "$ROLE_NAME" 2>&1 | grep -q 'NoSuchEntity'; then
    echo "Creating role..."
    cat > trust-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
    aws iam create-role --role-name "$ROLE_NAME" --assume-role-policy-document file://trust-policy.json > /dev/null
    aws iam attach-role-policy --role-name "$ROLE_NAME" --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
    rm trust-policy.json
    echo -e "${GREEN}Role created. Waiting 10s for propagation...${NC}"
    sleep 10
else
    echo "Role exists."
fi

ROLE_ARN=$(aws iam get-role --role-name "$ROLE_NAME" --query 'Role.Arn' --output text)

# 2. Deploy Functions
for file in "$LAMBDA_DIR"/*.js; do
    filename=$(basename -- "$file")
    func_name="${filename%.*}"
    zip_file="${func_name}.zip"

    echo -e "\n${BLUE}Processing $func_name...${NC}"

    # Zip the file
    zip -j "$zip_file" "$file" > /dev/null

    # Check if function exists
    if aws lambda get-function --function-name "$func_name" 2>&1 | grep -q 'ResourceNotFoundException'; then
        echo "Creating function..."
        aws lambda create-function \
            --function-name "$func_name" \
            --runtime "$RUNTIME" \
            --role "$ROLE_ARN" \
            --handler "${func_name}.handler" \
            --zip-file "fileb://${zip_file}" \
            --timeout 15 \
            > /dev/null
        echo -e "${GREEN}Created $func_name${NC}"
    else
        echo "Updating function code..."
        aws lambda update-function-code \
            --function-name "$func_name" \
            --zip-file "fileb://${zip_file}" \
            > /dev/null
        echo -e "${GREEN}Updated $func_name${NC}"
    fi

    # Cleanup zip
    rm "$zip_file"
done

echo -e "\n${GREEN}Deployment Complete!${NC}"
