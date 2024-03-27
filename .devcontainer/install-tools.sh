#!/bin/sh

GO_SWAGGER_VERSION="0.30.5"
TEMPL_VERSION="0.2.639"

ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
  GO_SWAGGER_BINARY_NAME="swagger_linux_amd64"
  TEMPL_BUNDLE_NAME="templ_Linux_x86_64.tar.gz"
elif [ "$ARCH" == "aarch64" ]; then
  GO_SWAGGER_BINARY_NAME="swagger_linux_arm64"
  TEMPL_BUNDLE_NAME="templ_Linux_arm64.tar.gz"
else
  echo "Unsupported architecture '${ARCH}'. Aborting."
  exit 1
fi

echo "Downloading Go-Swagger (${GO_SWAGGER_BINARY_NAME}) ..."
GO_SWAGGER_DL_URL="https://github.com/go-swagger/go-swagger/releases/download/v${GO_SWAGGER_VERSION}/${GO_SWAGGER_BINARY_NAME}"
sudo curl -L "${GO_SWAGGER_DL_URL}" -o /usr/local/bin/go-swagger
sudo chmod +x /usr/local/bin/go-swagger

echo "Downloading Templ (${TEMPL_BUNDLE_NAME}) ..."
TEMPL_DL_URL="https://github.com/a-h/templ/releases/download/v${TEMPL_VERSION}/${TEMPL_BUNDLE_NAME}"
curl -L "${TEMPL_DL_URL}" | sudo tar -xz -C /usr/local/bin templ
sudo chmod +x /usr/local/bin/templ
