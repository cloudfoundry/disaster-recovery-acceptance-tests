#!/usr/bin/env bash

set -eu

export CF_DEPLOYMENT_NAME
export CF_API_URL
export CF_ADMIN_USERNAME
export CF_ADMIN_PASSWORD
export BOSH_CERT_PATH
export BOSH_CLIENT
export BOSH_CLIENT_SECRET
export BOSH_ENVIRONMENT
export BBR_BUILD_PATH

go get github.com/onsi/ginkgo/ginkgo
glide install --strip-vendor
ginkgo -v -r --trace .
