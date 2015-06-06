default: deps

deps:
	go get github.com/tools/godep
	godep restore
