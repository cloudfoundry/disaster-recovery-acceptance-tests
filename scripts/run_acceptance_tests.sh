#!/usr/bin/env bash

set -eu

export DEPLOYMENT_TO_BACKUP=cf-integration-0
export DEPLOYMENT_TO_RESTORE=cf-integration-0
export BOSH_CERT_PATH=/Users/pivotal/workspace/bosh-backup-and-restore-meta/certs/genesis-bosh.backup-and-restore.cf-app.com.crt
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=DWTk6zUf6ogtPzl501e4Wi0eGNBQSS
export BOSH_URL=https://genesis-bosh.backup-and-restore.cf-app.com
export BOSH_GATEWAY_USER=vcap
export BOSH_GATEWAY_HOST=genesis-bosh.backup-and-restore.cf-app.com
export BOSH_GATEWAY_KEY=/Users/pivotal/workspace/bosh-backup-and-restore-meta/genesis-bosh/bosh.pem
export BBR_BUILD_PATH=/Users/pivotal/workspace/go/src/github.com/pivotal-cf/bosh-backup-and-restore/bbr

go get github.com/onsi/ginkgo/ginkgo
glide install
ginkgo -v -r --trace .
