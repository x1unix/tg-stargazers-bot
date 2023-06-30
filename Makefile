GO ?= go
WIRE ?= wire

DOCKERFILE=./Dockerfile
IMG_NAME ?= x1unix/tg-stargazers-bot
ENV_FILE ?= .env

.PHONY: run
run:
	@go run ./cmd/stargazers-server -f $(ENV_FILE)

.PHONY: wire
wire:
	@echo ":: Wiring dependencies..." && $(WIRE) gen ./internal/app

.PHONY:build
build:
	@if [ -z "$(TAG)" ]; then\
		echo "required parameter TAG is undefined" && exit 1; \
	fi;
	@echo ":: Building '$(IMG_NAME):latest' $(TAG)..." && \
	docker image build -t $(IMG_NAME):latest -t $(IMG_NAME):$(TAG) -f $(DOCKERFILE) \
		--build-arg APP_VERSION=$(TAG) .