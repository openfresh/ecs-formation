GOTEST_FLAGS=-cpu=1,2,4

default: deps

BASE_PACKAGE=github.com/stormcat24/ecs-formation
PACKAGES=util \
		aws \
		task

TEST_TARGETS=$(addprefix test-,$(PACKAGES))

deps:
	go get github.com/Masterminds/glide
	GO15VENDOREXPERIMENT=1 glide update --cache

deps-test:
	go get github.com/Masterminds/glide
	GO15VENDOREXPERIMENT=1 glide update --cache
	go get github.com/golang/lint/golint
	go get github.com/jstemmer/go-junit-report

build:
	GO15VENDOREXPERIMENT=1 go build -o bin/ecs-formation main.go

vet:
	GO15VENDOREXPERIMENT=1 go vet $(shell GO15VENDOREXPERIMENT=1 go list github.com/stormcat24/ecs-formation/... | grep -v vendor)

test: $(TEST_TARGETS)

$(TEST_TARGETS): test-%:
		@echo "**********************************************************"
		@echo " testing package: $*"
		@echo "**********************************************************"
		GO15VENDOREXPERIMENT=1 go test -v -covermode=atomic -coverprofile=coverage.out $(GOTEST_FLAGS) $(BASE_PACKAGE)/$(*)
		GO15VENDOREXPERIMENT=1 go test -v -run=nonthing -benchmem -bench=".*" $(GOTEST_FLAGS) $(BASE_PACKAGE)/$(*)

ci-test: $(CI_TEST_TARGETS)

$(CI_TEST_TARGETS): ci-test-%:
		@echo "**********************************************************"
		@echo " testing package: $*"
		@echo "**********************************************************"
		GO15VENDOREXPERIMENT=1 go test -covermode=atomic -coverprofile=coverage.out $(GOTEST_FLAGS) $(BASE_PACKAGE)/$(*)
		mkdir -p $(CIRCLE_ARTIFACTS)/$@/
		GO15VENDOREXPERIMENT=1 go tool cover -html=coverage.out -o $(CIRCLE_ARTIFACTS)/$@/cover.html
		mkdir -p $(CIRCLE_TEST_REPORTS)/junit
		cat coverage.out | go-junit-report > $(CIRCLE_TEST_REPORTS)/junit/$(notdir $@).xml
