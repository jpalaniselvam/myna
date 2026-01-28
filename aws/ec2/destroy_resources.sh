#!/bin/bash
set -e
INSTANCE_ID=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=myna-test-instance" "Name=instance-state-name,Values=pending,running,stopped,stopping" --query "Reservations[].Instances[].InstanceId" --output text)

if [ -n "$INSTANCE_ID" ]; then
    echo "Terminating Instance: $INSTANCE_ID"
    aws ec2 terminate-instances --instance-ids $INSTANCE_ID
    aws ec2 wait instance-terminated --instance-ids $INSTANCE_ID
else
    echo "Instance not found."
fi
