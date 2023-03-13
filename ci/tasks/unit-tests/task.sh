#!/usr/bin/env bash

set -eu

eval "$(ssh-agent)"
github_ssh_key_file=$(mktemp)

echo "$GITHUB_SSH_KEY" > "$github_ssh_key_file"
chmod 400 "$github_ssh_key_file"
ssh-add "$github_ssh_key_file"

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

pushd src/github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests
    ginkgo -v -r --trace --skipPackage acceptance
popd
