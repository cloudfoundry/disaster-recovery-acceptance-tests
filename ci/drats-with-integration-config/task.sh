#!/usr/bin/env bash

set -eu

export CONFIG="$PWD/drats-integration-config/${CONFIG_FILE_PATH}"

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

BOSH_GW_USER=$(jq -r .ssh_proxy_user "${CONFIG}")
BOSH_GW_HOST=$(jq -r .ssh_proxy_host "${CONFIG}")
BOSH_GW_PRIVATE_KEY=$(jq -r .ssh_proxy_private_key "${CONFIG}")

rm -f ~/.gitconfig
echo "${BOSH_GW_PRIVATE_KEY}" > ssh.key
chmod 0400 ssh.key

BOSH_ALL_PROXY="ssh+socks5://${BOSH_GW_USER}@${BOSH_GW_HOST}:22?private-key=${PWD}/ssh.key"
export BOSH_ALL_PROXY

pushd bbr-binary-release
  tar xvf ./*.tar
  export BBR_BUILD_PATH="$PWD/releases/bbr"
popd

pushd src/github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
