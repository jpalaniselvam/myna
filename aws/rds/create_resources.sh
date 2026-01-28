#!/bin/bash
set -e
DB_ID="myna-test-db"

STATUS=$(aws rds describe-db-instances --db-instance-identifier $DB_ID --query "DBInstances[0].DBInstanceStatus" --output text 2>/dev/null || echo "")

if [ -z "$STATUS" ]; then
    echo "Creating DB Instance..."
    aws rds create-db-instance \
        --db-instance-identifier $DB_ID \
        --db-instance-class db.t3.micro \
        --engine postgres \
        --master-username postgres \
        --master-user-password mypassword123 \
        --allocated-storage 20 \
        --publicly-accessible \
        --tags "Key=Project,Value=myna-test"
    echo "Creation initiated."
else
    echo "DB Instance exists: $STATUS"
fi
