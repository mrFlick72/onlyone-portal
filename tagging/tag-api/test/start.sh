#! /bin/bash

export AWS_ACCESS_KEY_ID="xxx"
export AWS_SECRET_ACCESS_KEY="xxx"
export AWS_REGION="us-east-1"
export DYNAMODB_ENDPOINT="http://localhost:4566"
export TAGS_TABLE_NAME="Tags_Local_Test_Table"
export OAUTH2_ROLE="ROLE_USER"
export JWKS_ENDPOINT="http://local.api.vauthenticator.com:9090/oauth2/jwks"
export CORS_ALLOWED_ORIGINS="http://local.onlyone-portal.com:8070"

go run ../main.go
