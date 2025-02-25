.PHONY: mod build
.DEFAULT_GOAL := build

BUILD = `date +%FT%T%z`

mod:
	go mod verify

build:
	go build -ldflags "-w -s -X main.Build=${BUILD}" ${GOPATH}/src/github.com/arxon31/sso/cmd/sso/sso.go

tests:
	go install github.com/matryer/moq@v0.3.4
	go generate ./...
	go test -v -p 1 ./...

coverage:
	go install github.com/matryer/moq@v0.3.4
	go generate ./...
	go test -coverprofile=coverage.out -v -p 1 ./...

lint:
	go install github.com/matryer/moq@v0.3.4
	go generate ./...
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./... --config=./.golangci.yml

format:
	goimports -local "github.com/arxon31/sso" -w `pwd`

run:
	docker-compose -f ./deployments/docker-compose-dev.yml up -d --build