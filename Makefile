GOTEST_FLAGS=-cpu=1,2,4

default: deps

BASE_PACKAGE=github.com/stormcat24/ecs-formation
PACKAGES=util \
		aws

TEST_TARGETS=$(addprefix test-,$(PACKAGES))


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

test: $(TEST_TARGETS)

$(TEST_TARGETS): test-%:
		@echo "**********************************************************"
		@echo " testing package: $*"
		@echo "**********************************************************"
		GO15VENDOREXPERIMENT=1 go test -v -covermode=atomic -coverprofile=coverage.out $(GOTEST_FLAGS) $(BASE_PACKAGE)/$(*)
		GO15VENDOREXPERIMENT=1 go test -v -run=nonthing -benchmem -bench=".*" $(GOTEST_FLAGS) $(BASE_PACKAGE)/$(*)

