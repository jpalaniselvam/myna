#!/bin/bash
set -e
echo "Fetching AMI..."
AMI_ID=$(aws ssm get-parameters --names /aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64 --query "Parameters[0].Value" --output text)
echo "Using AMI: $AMI_ID"

INSTANCE_ID=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=myna-test-instance" "Name=instance-state-name,Values=pending,running,stopped,stopping" --query "Reservations[].Instances[].InstanceId" --output text)

if [ -z "$INSTANCE_ID" ]; then
    echo "Creating Instance..."
    ID=$(aws ec2 run-instances --image-id $AMI_ID --instance-type t3.micro --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=myna-test-instance}]" --query "Instances[0].InstanceId" --output text)
    echo "Created Instance: $ID"
    aws ec2 wait instance-running --instance-ids $ID
else
    echo "Instance exists: $INSTANCE_ID"
fi
