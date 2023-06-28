#!/usr/bin/env bash
set -e

KEY_NAME="jwt"
DEST_DIR="$PWD"

die() { echo "$*" >&2; exit 2; }
needs_arg() { if [ -z "$OPTARG" ]; then die "Missing required argument for --$OPT option"; fi; }

show_help() {
  cat << EOF
gen-keys.sh - Generates JWT public and private keys.

Usage: gen-keys.sh [-n] [-d] [-h|--help]

Options:
  -n                Set key name          (default: "$KEY_NAME")
  -d                Set output directory  (default: "$DEST_DIR")
  -h, --help        Show help

EOF
exit
}

while getopts "h:n:d:-:" OPT; do
  # support long options: https://stackoverflow.com/a/28466267/519360
  if [ "$OPT" = "-" ]; then   # long option: reformulate OPT and OPTARG
    OPT="${OPTARG%%=*}"       # extract long option name
    OPTARG="${OPTARG#$OPT}"   # extract long option argument (may be empty)
    OPTARG="${OPTARG#=}"      # if long option argument, remove assigning `=`
  fi
  case "$OPT" in
    h|help)
      show_help
      exit
      ;;
    n) needs_arg; KEY_NAME="$OPTARG" ;;
    d) needs_arg; DEST_DIR="$OPTARG" ;;
    ??* ) die "Illegal option --$OPT" ;;
    ? ) exit 2 ;;
  esac
done

if [ ! -d "$DEST_DIR" ]; then
  die "ERROR: Output directory doesn't exists: '$DEST_DIR'"
fi

if ! command -v "openssl" &> /dev/null
then
    die "ERROR: OpenSSL not found. Please install openssl package."
fi

key_file="$DEST_DIR/$KEY_NAME.key"
pub_file="$DEST_DIR/$KEY_NAME.pub"

echo ":: Generating private key '$key_file' ..."
openssl genrsa -out "$key_file" 2048
echo ":: Generating public key '$pub_file' ..."
openssl rsa -in "$key_file" -outform PEM -pubout -out "$pub_file"
cat << EOF
:: Generated JWT keypair:
Private key: "$key_file"
Public key:  "$pub_file"
EOF