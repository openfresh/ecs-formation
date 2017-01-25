GOTEST_FLAGS=-cpu=1,2,4

default: deps

BASE_PACKAGE=github.com/stormcat24/ecs-formation

IGNORE=vendor|cache$$
TARGETS=$(shell go list ./... | awk '$$0 !~ /$(IGNORE)/{print $0}')
ARCH=$(shell uname | tr '[:upper:]' '[:lower:]')

deps:
		go get github.com/golang/dep
		go get github.com/golang/lint/golint
		go get github.com/jstemmer/go-junit-report
		go get github.com/golang/mock/gomock
		go get github.com/golang/mock/mockgen
		dep ensure

build:
		go build -o bin/ecs-formation main.go

test: vet
		go test -cover $(TARGETS)

vet:
		go vet $(TARGETS)

mock:
	  mockgen -source client/applicationautoscaling/client.go -package applicationautoscaling -destination client/applicationautoscaling/mock.go
	  mockgen -source client/autoscaling/client.go -package autoscaling -destination client/autoscaling/mock.go
	  mockgen -source client/elb/client.go -package elb -destination client/elb/mock.go
	  mockgen -source client/elbv2/client.go -package elbv2 -destination client/elbv2/mock.go
	  mockgen -source client/iam/client.go -package iam -destination client/iam/mock.go
	  mockgen -source client/s3/client.go -package s3 -destination client/s3/mock.go
