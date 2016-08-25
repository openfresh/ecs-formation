GOTEST_FLAGS=-cpu=1,2,4

default: deps

BASE_PACKAGE=github.com/stormcat24/ecs-formation

IGNORE=vendor|cache$$
TARGETS=$(shell go list ./... | awk '$$0 !~ /$(IGNORE)/{print $0}')
ARCH=$(shell uname | tr '[:upper:]' '[:lower:]')

deps:
	go get github.com/Masterminds/glide
	glide update --cache

deps-test:
	go get github.com/Masterminds/glide
	glide update --cache
	go get github.com/golang/lint/golint
	go get github.com/jstemmer/go-junit-report

update:
	rm -rf ./vendor
	glide update --cache

build:
	go build -o bin/ecs-formation main.go

test: vet
		go test -cover $(TARGETS)

vet:
		go vet $(TARGETS)