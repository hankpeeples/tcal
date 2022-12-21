BINARY=tcal

.PHONY: all build run clean test

all: build run clean

build: 
	@echo "--> mod tidy" && go mod tidy
	@echo "--> Building..." && GOARCH=amd64 GOOS=darwin go build -o tcal-bin github.com/hankpeeples/tcal

run: 
	@echo "--> Running" && ./tcal-bin

clean: 
	@echo "--> Cleaning" && go clean
	@rm tcal-bin
