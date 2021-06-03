#!/bin/bash

ARCH="amd64"
OWNER="$1" # change me
OUT_FILE="jaws-$(echo "${OWNER}" | tr '[:upper:]' '[:lower:]').exe"

GOOS=windows GOARCH=${ARCH} \
  go build \
  -ldflags "-X github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess.Owner=Theo" \
  -o "${OUT_FILE}" \
  main.go