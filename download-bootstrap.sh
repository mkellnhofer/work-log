#!/bin/sh

BOOTSTRAP_VERSION="5.3.2"
BOOTSTRAP_BUNDLE_NAME="bootstrap-${BOOTSTRAP_VERSION}.tar.gz"
BOOTSTRAP_DL_URL="https://github.com/twbs/bootstrap/archive/refs/tags/v${BOOTSTRAP_VERSION}.tar.gz"

rm -rf /workspace/scss/vendor/bootstrap

echo "Downloading Bootstrap ${BOOTSTRAP_VERSION} ..."

cd /tmp
curl -sL ${BOOTSTRAP_DL_URL} -o ${BOOTSTRAP_BUNDLE_NAME}
tar xzf ${BOOTSTRAP_BUNDLE_NAME} bootstrap-${BOOTSTRAP_VERSION}/scss/
cp -r bootstrap-${BOOTSTRAP_VERSION}/scss /workspace/scss/vendor/bootstrap
rm -rf bootstrap-${BOOTSTRAP_VERSION}
rm -f ${BOOTSTRAP_BUNDLE_NAME}
