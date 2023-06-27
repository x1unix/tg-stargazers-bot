GO ?= go
WIRE ?= wire

ENV_FILE ?= .env

.PHONY: run
run:
	@go run ./cmd/stargazers-server -f $(ENV_FILE)

.PHONY: wire
wire:
	@echo ":: Wiring dependencies..." && $(WIRE) gen ./internal/app
