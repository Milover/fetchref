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
	go build -ldflags "$(LDFLAGS)" -o $(TARGET) main.go

run:
	./$(TARGET)

test:
	go test ./...

vet:
	go vet ./...

lint:
	$(shell go env GOPATH)/bin/golangci-lint run ./...

clean:
	go clean
	rm -f $(TARGET)
