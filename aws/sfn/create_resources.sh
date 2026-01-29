#!/bin/bash
set -e
ROLE="myna-test-sfn-role"
SM="myna-test-state-machine"

ROLE_ARN=$(aws iam get-role --role-name $ROLE --query "Role.Arn" --output text 2>/dev/null || echo "")
if [ -z "$ROLE_ARN" ]; then
    echo "Creating Role..."
    cat > trust.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "states.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
    ROLE_ARN=$(aws iam create-role --role-name $ROLE --assume-role-policy-document file://trust.json --query "Role.Arn" --output text)
    rm trust.json
fi

SM_ARN=$(aws stepfunctions list-state-machines --query "stateMachines[?name=='$SM'].stateMachineArn" --output text)
if [ -z "$SM_ARN" ]; then
    echo "Creating State Machines..."
    aws stepfunctions create-state-machine --name $SM --role-arn $ROLE_ARN --definition '{"StartAt": "Pass", "States": {"Pass": {"Type": "Pass", "End": true}}}'
else
    echo "State Machine exists."
fi
