#!/usr/bin/env bash

set -eu

export BOSH_CERT_PATH
export BOSH_GATEWAY_KEY
export BOSH_CLIENT
export BOSH_URL
export BOSH_GATEWAY_USER
export BOSH_GATEWAY_HOST
export BBR_BUILD_PATH

pushd src/github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests
    go get github.com/onsi/ginkgo/ginkgo
    glide install
	ginkgo -v -r --trace acceptance
  popd
popd
