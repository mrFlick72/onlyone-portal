
go mod edit -replace github.com/mrflick72/onlyone-portal/core-services/golang-web-framework=../../core-services/golang-web-framework

replace github.com/mrflick72/onlyone-portal/core-services/golang-web-framework => ../../core-services/golang-web-framework

go mod tidy


docker build -t mrflick72/budget/budget-api:1 -f ../../core-services/docker/ubuntu.Dockerfile  .