#!/bin/bash

echo "Initializing DynamoDB tables..."

# Wait for LocalStack to be ready
until awslocal dynamodb list-tables; do
  echo "Waiting for DynamoDB to be ready..."
  sleep 1
done

# Create accounts table
awslocal dynamodb create-table \
  --table-name accounts \
  --attribute-definitions \
    AttributeName=id,AttributeType=N \
  --key-schema \
    AttributeName=id,KeyType=HASH \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5

echo "DynamoDB table 'accounts' created successfully"

awslocal dynamodb create-table \
  --table-name transactions \
  --attribute-definitions \
    AttributeName=id,AttributeType=S \
  --key-schema \
    AttributeName=id,KeyType=HASH \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5

echo "DynamoDB table 'accounts' created successfully"

# Insert initial data
awslocal dynamodb put-item \
  --table-name accounts \
  --item '{
    "id": {"N": "1"},
    "name": {"S": "main"},
    "created_at": {"S": "2025-01-23 12:46:20"}
  }'

echo "Initial data inserted successfully"

