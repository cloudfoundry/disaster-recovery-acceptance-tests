#!/usr/bin/env bash

set -eu

eval "$(ssh-agent)"
github_ssh_key=$(mktemp)
echo "$GITHUB_SSH_KEY" > "$github_ssh_key"
chmod 400 "$github_ssh_key"
ssh-add "$github_ssh_key"

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

pushd src/github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests
    go get github.com/onsi/ginkgo/ginkgo
    dep ensure
    ginkgo -v -r --trace --skipPackage acceptance
popd