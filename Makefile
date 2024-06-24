.PHONY: check view-coverage

GOFILES := $(shell find ./internal ./cmd -name '*.go')

all: check bin barny

barny: $(GOFILES) bin
	go build -o bin/ ./cmd/barny

bin:
	mkdir bin

check:
	go fmt ./...
	go vet ./...
	staticcheck ./...
	go test ./... -race -covermode=atomic -coverprofile=coverage.txt -shuffle on

view-coverage:
	go tool cover -html=coverage.txt
