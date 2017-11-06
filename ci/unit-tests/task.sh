#!/usr/bin/env bash

set -eu

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

pushd src/github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests
    go get github.com/onsi/ginkgo/ginkgo
    dep ensure
    ginkgo -v -r --trace --skipPackage acceptance
popd