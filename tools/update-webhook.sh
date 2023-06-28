#!/usr/bin/env bash
set -e

die() { echo "$*" >&2; exit 2; }
needs_arg() { if [ -z "$OPTARG" ]; then die "Missing required argument for --$OPT option"; fi; }

urlencode() {
  perl -MURI::Escape -e 'print uri_escape($ARGV[0]);' "$1"
}

joinpath() {
  local basepath=$1
  local relpath=$2

  # Remove trailing slash from the base path, if present
  basepath=${basepath%/}

  # Remove leading slash from the relative path, if present
  relpath=${relpath#/}

  # Join the base path and relative path with a forward slash
  local joined_path="$basepath/$relpath"

  echo "$joined_path"
}

show_help() {
  cat << EOF
update-webhook.sh - Update Telegram webhook URL

Usage: update-webhook.sh env_file_path [-h|--help]

Options:
  -h, --help        Show help

EOF
exit
}

if [ -z "$1" ]; then
  echo "Error: missing env file argument"
  show_help
  exit
fi

case "$1" in
  "-h"|"--help")
    show_help
    exit
    ;;
  *) ;;
esac

if ! command -v perl >/dev/null 2>&1; then
  die "Error: Perl is required"
fi

env_file="$1"
if [ ! -f "$env_file" ]; then
  die "Error: env file $env_file doesn't exists"
fi

# shellcheck disable=SC1090
. "$env_file"

keys=('HTTP_BASE_URL' 'TELEGRAM_BOT_TOKEN' 'BOT_WEBHOOK_SECRET')
for key in "${keys[@]}"; do
  if [ -z "${!key}" ]; then
    die "Error: missing $key in environment file"
  fi
done

webhook_url="/webhook/telegram?s=$(urlencode "$BOT_WEBHOOK_SECRET")"
webhook_url=$(joinpath "$HTTP_BASE_URL" "$webhook_url")

echo ":: Updating webhook to '$webhook_url'..."

curl -X POST -w "\n" \
  --fail-with-body \
  "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook" \
  --data-urlencode "url=$webhook_url"

echo ":: New Webhook info:"
curl -w "\n" -f "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getWebhookInfo"
