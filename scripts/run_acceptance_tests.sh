#!/usr/bin/env bash

set -eu

export DEPLOYMENT_TO_BACKUP
export DEPLOYMENT_TO_RESTORE
export BOSH_CERT_PATH
export BOSH_CLIENT
export BOSH_CLIENT_SECRET
export BOSH_URL
export BBR_BUILD_PATH

go get github.com/onsi/ginkgo/ginkgo
glide install --strip-vendor
ginkgo -v -r --trace .
