#!/usr/bin/env bash

set -eu

export CONFIG="$PWD/drats-integration-config/${CONFIG_FILE_PATH}"

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

BOSH_GW_USER=$(jq -r .ssh_proxy_user "${CONFIG}")
BOSH_GW_HOST=$(jq -r .ssh_proxy_host "${CONFIG}")
BOSH_GW_PRIVATE_KEY=$(jq -r .ssh_proxy_private_key "${CONFIG}")
SSH_DESTINATION_CIDR=$(jq -r .ssh_proxy_cidr "${CONFIG}")

eval "$(ssh-agent)"

rm -f ~/.gitconfig
echo "${BOSH_GW_PRIVATE_KEY}" > ssh.pem
chmod 0400 ssh.pem
ssh-add ssh.pem

sshuttle -r "${BOSH_GW_USER}@${BOSH_GW_HOST}" "${SSH_DESTINATION_CIDR}" --daemon -e 'ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o ServerAliveInterval=600' --no-latency-control

sleep 5

if ! stat sshuttle.pid > /dev/null 2>&1; then
  echo "Failed to start sshuttle daemon"
  exit 1
fi

pushd bbr-binary-release
  tar xvf ./*.tar
  export BBR_BUILD_PATH="$PWD/releases/bbr"
popd

pushd src/github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
