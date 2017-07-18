#!/usr/bin/env bash

set -eu

#eval "$(ssh-agent)"
#./bosh-backup-and-restore-meta/unlock-ci.sh
#chmod 400 bosh-backup-and-restore-meta/keys/github
#chmod 400 bosh-backup-and-restore-meta/genesis-bosh/bosh.pem
#ssh-add bosh-backup-and-restore-meta/keys/github

#export META_PATH=`pwd`/bosh-backup-and-restore-meta
export META_PATH=/Users/pivotal/workspace/bosh-backup-and-restore-meta

#export GOPATH=$PWD
#export PATH=$PATH:$GOPATH/bin
export BOSH_CERT_PATH="${META_PATH}"/certs/genesis-bosh.backup-and-restore.cf-app.com.crt
export BOSH_CLIENT=admin
export BOSH_URL=https://genesis-bosh.backup-and-restore.cf-app.com
export BOSH_GATEWAY_USER=vcap
export BOSH_GATEWAY_HOST=genesis-bosh.backup-and-restore.cf-app.com
export BOSH_GATEWAY_KEY="${META_PATH}"/genesis-bosh/bosh.pem


export BBR_BUILD_PATH=$PWD/bbr-binary-release/

#pushd src/github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests
#    go get github.com/onsi/ginkgo/ginkgo

    glide install

	ginkgo -v -r --trace acceptance
#popd
