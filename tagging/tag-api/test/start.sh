#! bin/bash

export AWS_ACCESS_KEY_ID="xxx"
export AWS_SECRET_ACCESS_KEY="xxx"
export AWS_REGION="us-east-1"
export DYNAMODB_ENDPOINT="http://localhost:4566"

go run main.go
go run ../main.go
