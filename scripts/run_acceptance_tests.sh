#!/usr/bin/env bash

set -e

if [ -z "$DEPLOYMENT_TO_BACKUP" ]; then
  echo "Need to set DEPLOYMENT_TO_BACKUP"
  exit 1
fi

if [ -z "$DEPLOYMENT_TO_RESTORE" ]; then
  echo "Need to set DEPLOYMENT_TO_RESTORE"
  exit 1
fi

export BOSH_META_PATH=~/workspace/bosh-backup-and-restore-meta

export TEAM_GPG_KEY="$(lpass show "Shared-PCF-Backup-and-Restore/CF Lazarus GPG key" --notes)"
export BOSH_CLIENT_SECRET="$(lpass show "GenesisBoshDirectorGCP" --password)"
export BOSH_CERT_PATH="$BOSH_META_PATH"/certs/genesis-bosh.backup-and-restore.cf-app.com.crt
export BOSH_CLIENT=admin
export BOSH_URL=https://genesis-bosh.backup-and-restore.cf-app.com
export BOSH_GATEWAY_USER=vcap
export BOSH_GATEWAY_HOST=genesis-bosh.backup-and-restore.cf-app.com
export BOSH_GATEWAY_KEY="$BOSH_META_PATH"/genesis-bosh/bosh.pem


pushd "$GOPATH"/src/github.com/pivotal-cf/bosh-backup-and-restore
  make bin-linux
  export BBR_BUILD_PATH=$PWD/bbr
popd

pushd "$(dirname $0)"; pushd ..
    glide install
	ginkgo -v -r --trace acceptance
popd; popd

rm "$BBR_BUILD_PATH"