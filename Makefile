.PHONY: build install test clean

BINARY := repowiki
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/repowiki

install:
	go install ./cmd/repowiki

test:
	go test ./internal/... -v -race

clean:
	rm -rf $(BUILD_DIR)
