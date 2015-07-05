GOTEST_FLAGS=-cpu=1,2,4

default: deps

BASE_PACKAGE=github.com/stormcat24/ecs-formation
PACKAGES=util

TEST_TARGETS=$(addprefix test-,$(PACKAGES))


deps:
	go get github.com/tools/godep
	godep restore

deps-save:
	godep save $(BASE_PACKAGE)/...

deps-test:
	go get github.com/tools/godep
	godep restore
	go get github.com/golang/lint/golint
	go get github.com/jstemmer/go-junit-report

build:
	godep go build -o bin/ecs-formation main.go

test: $(TEST_TARGETS)

$(TEST_TARGETS): test-%:
		@echo "**********************************************************"
		@echo " testing package: $*"
		@echo "**********************************************************"
		cd $* && godep go test -v -covermode=atomic -coverprofile=coverage.out $(GOTEST_FLAGS) | tee test.out && test $${PIPESTATUS[0]} -eq 0
		cd $* && cat test.out | go-junit-report > test.xml
		cd $* && godep go tool cover -html=coverage.out -o coverage.html


