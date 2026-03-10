.PHONY: help build run test clean install

BINARY := cs
CMD_DIR := ./cmd/cs
LOCAL_BIN := $(HOME)/.local/bin

help:
	@printf "Available targets:\n"
	@printf "  make build       Build the cs binary\n"
	@printf "  make run ARGS='' Run cs with arguments\n"
	@printf "  make test        Run all Go tests\n"
	@printf "  make clean       Remove built binary\n"
	@printf "  make install     Install binary to ~/.local/bin\n"

build:
	go build -o $(BINARY) $(CMD_DIR)

run: build
	./$(BINARY) $(ARGS)

test:
	go test ./...

clean:
	rm -f ./$(BINARY)

install: build
	mkdir -p "$(LOCAL_BIN)"
	cp ./$(BINARY) "$(LOCAL_BIN)/$(BINARY)"
	chmod +x "$(LOCAL_BIN)/$(BINARY)"
