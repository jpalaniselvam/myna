#!/bin/bash
set -e
SM="myna-test-state-machine"
ROLE="myna-test-sfn-role"

SM_ARN=$(aws stepfunctions list-state-machines --query "stateMachines[?name=='$SM'].stateMachineArn" --output text)
if [ -n "$SM_ARN" ]; then
    echo "Deleting State Machine..."
    aws stepfunctions delete-state-machine --state-machine-arn $SM_ARN
fi

echo "Deleting Role..."
aws iam delete-role --role-name $ROLE || echo "Role not found."
