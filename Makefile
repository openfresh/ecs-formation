GOTEST_FLAGS=-cpu=1,2,4

default: test

deps:
	go get github.com/tools/godep
	godep restore

test: generate
	TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4
	@$(MAKE) vet
