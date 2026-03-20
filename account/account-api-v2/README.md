# account-service-ui

go mod edit -replace github.com/mrflick72/onlyone-portal/core-services/golang-web-framework=../../core-services/golang-web-framework

replace github.com/mrflick72/onlyone-portal/core-services/golang-web-framework => ../../core-services/golang-web-framework

go mod tidy

CGO_ENABLED=0 GOOS=linux go build -o app .

docker build -t mrflick72/account/account-api-v2:1 -f ../../core-services/docker/ubuntu.Dockerfile  .