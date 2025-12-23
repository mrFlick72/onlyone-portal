#! bin/bash

export AWS_ACCESS_KEY_ID="xxx"
export AWS_SECRET_ACCESS_KEY="xxx"
export AWS_REGION="us-east-1"
export DYNAMODB_ENDPOINT="http://localhost:4566"
export OAUTH2_ROLE="ROLE_USER"
export JWKS_ENDPOINT="http://local.api.vauthenticator.com:9090/oauth2/jwks"

go run main.go
go run ../main.go
