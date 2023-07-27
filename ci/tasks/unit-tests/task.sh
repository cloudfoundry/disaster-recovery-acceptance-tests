#!/usr/bin/env bash

set -eu

eval "$(ssh-agent)"
github_ssh_key_file=$(mktemp)

echo "$GITHUB_SSH_KEY" > "$github_ssh_key_file"
chmod 400 "$github_ssh_key_file"
ssh-add "$github_ssh_key_file"

pushd src/github.com/cloudfoundry/disaster-recovery-acceptance-tests
    go run github.com/onsi/ginkgo/v2/ginkgo -v -r --trace --skip-package acceptance
popd
