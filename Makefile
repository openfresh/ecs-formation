GOTEST_FLAGS=-cpu=1,2,4

default: deps

deps:
	go get github.com/Masterminds/glide
	GO15VENDOREXPERIMENT=1 glide update

deps-save:
	glide update

deps-test:
	go get github.com/Masterminds/glide
	GO15VENDOREXPERIMENT=1 glide update
	go get github.com/golang/lint/golint
	go get github.com/jstemmer/go-junit-report

build:
	GO15VENDOREXPERIMENT=1 go build -o bin/ecs-formation main.go

test:
	GO15VENDOREXPERIMENT=1 go test $(go list ./... | grep -v vendor)
