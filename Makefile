# Makefile

TARGET		:= pressurized-cylinder	

build:
	go build -o $(TARGET) main.go

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
