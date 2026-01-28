#!/bin/bash
set -e
DB_ID="myna-test-db"
echo "Deleting DB Instance..."
aws rds delete-db-instance --db-instance-identifier $DB_ID --skip-final-snapshot --delete-automated-backups || echo "Not found."
