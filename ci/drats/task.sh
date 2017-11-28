#!/usr/bin/env bash

set -eu

export CF_DEPLOYMENT_NAME
export CF_ADMIN_USERNAME
export CF_ADMIN_PASSWORD
export CF_API_URL
export BOSH_CLIENT
export BOSH_CLIENT_SECRET
export BOSH_CA_CERT
export BOSH_ENVIRONMENT
export NFS_SERVICE_NAME
export NFS_PLAN_NAME
export NFS_BROKER_USER
export NFS_BROKER_PASSWORD
export NFS_BROKER_URL
export DELETE_AND_REDEPLOY_CF

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

eval "$(ssh-agent)"

rm -f ~/.gitconfig
echo "${BOSH_GW_PRIVATE_KEY}" > ssh.pem
chmod 0400 ssh.pem
ssh-add ssh.pem

sshuttle -r "${BOSH_GW_USER}@${BOSH_GW_HOST}" "${SSH_DESTINATION_CIDR}" --daemon -e 'ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o ServerAliveInterval=600'

sleep 5

echo "$BOSH_CA_CERT" > bosh.cert
export BOSH_CERT_PATH=`pwd`/bosh.cert

pushd bbr-binary-release
  tar xvf *.tar
  export BBR_BUILD_PATH=`pwd`/releases/bbr
popd

pushd src/github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
