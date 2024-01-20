#!/bin/sh

ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
  ARCH_ALT="amd64"
elif [ "$ARCH" == "aarch64" ]; then
  ARCH_ALT="arm64"
elif [ "$ARCH" == "armv7l" ]; then
  ARCH_ALT="arm"
else
  echo "Unsupported architecture '${ARCH}'. Aborting."
  exit 1
fi

GO_SWAGGER_VERSION="0.30.5"
GO_SWAGGER_BINARY="swagger_linux_${ARCH_ALT}"
GO_SWAGGER_DL_URL="https://github.com/go-swagger/go-swagger/releases/download/v${GO_SWAGGER_VERSION}/${GO_SWAGGER_BINARY}"

echo "Downloading Go-Swagger (${ARCH_ALT}) ..."

curl -L "${GO_SWAGGER_DL_URL}" -o ./go-swagger
chmod +x ./go-swagger
