# Makefile

MODULE		:= $(shell go list -m)
URL			:= https://$(MODULE)
TARGET		:= $(shell basename $(MODULE))
VERSION		:= $(shell git describe --tags --abbrev=0)
MAINTAINER	:= milovic.ph@gmail.com

META_PREFIX	:= $(MODULE)/internal/metainfo
LDFLAGS		:= -X '$(META_PREFIX).Project=$(TARGET)' \
			   -X '$(META_PREFIX).Version=$(VERSION)' \
			   -X '$(META_PREFIX).Url=$(URL)' \
			   -X '$(META_PREFIX).Maintainer=$(MAINTAINER)'

build:
	echo $(MODULE)
	go build -ldflags "$(LDFLAGS)" -o bin/$(TARGET) main.go

run:
	./...

test:
	go test ./...

testv:
	go test -v ./...

test-integration:
	go test -tags=integration ./...

vet:
	go vet ./...

lint:
	$(shell go env GOPATH)/bin/golangci-lint run ./...

update-deps:
	go get -u ./...
	go mod tidy

update-go:
	go mod edit -go=$(shell go version | awk '{print $$3}' | sed -e 's/go//g')
	go mod tidy

clean:
	go clean
	rm -rf bin

.PHONY: run test testv test-integration vet lint clean update-deps update-go
