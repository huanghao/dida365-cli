BINARY := dida365-cli
CMD := ./cmd/dida

.PHONY: build run test tidy fmt

build:
	go build -o $(BINARY) $(CMD)

run:
	go run $(CMD) --help

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w cmd internal
