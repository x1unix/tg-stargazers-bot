GO ?= go

ENV_FILE ?= .env

.PHONY: run
run:
	@go run ./cmd/stargazers-server -f $(ENV_FILE)
