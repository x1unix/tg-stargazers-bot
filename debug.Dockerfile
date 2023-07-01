FROM golang:1.19-alpine as builder
ARG APP_VERSION='1.0.0'

COPY . /opt/stargazers-bot
WORKDIR /opt/stargazers-bot

RUN apk add --no-cache  \
    tzdata \
    ca-certificates \
    curl \
    openssl \
    bash \
    alpine-sdk \
    delve

RUN go build \
    -o stargazers-server \
    -gcflags="all=-N -l" \
    -ldflags="-X 'github.com/x1unix/tg-stargazers-bot/internal/app.Version=$APP_VERSION-debug'" \
    ./cmd/stargazers-server

ENV PATH "$PATH:/opt/stargazers-bot"
ENV APP_ENV dev
ENV HTTP_LISTEN ':8080'
ENV HTTP_STATIC_DIR '/opt/stargazers-bot/public'

HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/ || exit 1

EXPOSE 8080
EXPOSE 2345

ENTRYPOINT dlv \
  --listen=:2345 \
  --headless=true \
  --api-version=2 \
  --accept-multiclient \
  exec ./stargazers-server
