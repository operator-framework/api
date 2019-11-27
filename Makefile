all: tidy fmt

tidy:
	go mod tidy

test:
	go test -v ./...

fmt:
	go fmt ./...
