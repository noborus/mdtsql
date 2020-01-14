BINARY_NAME := mdtsql
SRCS := $(shell git ls-files '*.go')

all: build

test: $(SRCS)
	go test ./...

build: $(BINARY_NAME)

$(BINARY_NAME): $(SRCS)
	go build -o $(BINARY_NAME) ./cmd/mdtsql

install:
	go install ./cmd/mdtsql

clean:
	rm -f $(BINARY_NAME)

.PHONY: all test build install clean
