#! /bin/bash

export OAUTH2_ROLE="ROLE_USER"
export JWKS_ENDPOINT="http://local.api.vauthenticator.com:9090/oauth2/jwks"
export CORS_ALLOWED_ORIGINS="http://local.onlyone-portal.com:8070"
export WEBSERVER_PORT="5055"

go run ../main.go
