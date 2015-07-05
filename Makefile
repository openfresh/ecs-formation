default: deps

BASE_PACKAGE=github.com/stormcat24/ecs-formation

deps:
	go get github.com/tools/godep
	godep restore

deps-save:
	godep save $(BASE_PACKAGE)/...

build:
	godep go build -o bin/ecs-formation main.go
