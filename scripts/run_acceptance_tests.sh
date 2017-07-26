#!/usr/bin/env bash

set -eu

export DEPLOYMENT_TO_BACKUP
export DEPLOYMENT_TO_RESTORE
export BOSH_CERT_PATH
export BOSH_CLIENT
export BOSH_CLIENT_SECRET
export BOSH_URL
export BOSH_GATEWAY_USER
export BOSH_GATEWAY_HOST
export BOSH_GATEWAY_KEY
export BBR_BUILD_PATH

go get github.com/onsi/ginkgo/ginkgo
glide install
ginkgo -v -r --trace .
