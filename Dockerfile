FROM golang:1.19-alpine as builder
ARG APP_VERSION=1.0.0
COPY . /tmp/work
WORKDIR /tmp/work
RUN go build \
    -o stargazers-server \
    -ldflags="-X 'github.com/x1unix/tg-stargazers-bot/internal/app.Version=$APP_VERSION'" \
    ./cmd/stargazers-server

FROM alpine:3.18
RUN apk add --no-cache  \
    tzdata \
    ca-certificates \
    curl \
    openssl \
    bash && \
    mkdir -p /opt/stargazers-bot

COPY --from=builder /tmp/work/stargazers-server /opt/stargazers-bot
COPY --from=builder /tmp/work/tools/gen-keys.sh /opt/stargazers-bot
COPY --from=builder /tmp/work/tools/update-webhook.sh /opt/stargazers-bot
COPY --from=builder /tmp/work/public /opt/stargazers-bot/public

ENV PATH "$PATH:/opt/stargazers-bot"
ENV APP_ENV prod
ENV HTTP_LISTEN ":8080"
ENV HTTP_STATIC_DIR '/opt/stargazers-bot/public'

HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/ || exit 1

EXPOSE 8080

ENTRYPOINT /opt/stargazers-bot/stargazers-server
