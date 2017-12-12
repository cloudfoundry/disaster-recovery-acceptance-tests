#!/bin/bash

set -eu -o pipefail

# ENV
: "${GOPATH:?}"
: "${CONFIG:?}"
# The following params are optional
: "${SKIP_SUITE_NAME:=""}"

tmpdir="$( mktemp -d /tmp/run-drats.XXXXXXXXXX )"

BOSH_GW_USER=`jq -r .ssh_proxy_user ${CONFIG}`
BOSH_GW_HOST=`jq -r .ssh_proxy_host ${CONFIG}`
BOSH_GW_PRIVATE_KEY=`jq -r .ssh_proxy_private_key ${CONFIG}`
SSH_DESTINATION_CIDR=`jq -r .ssh_proxy_cidr ${CONFIG}`

ssh_key="${tmpdir}/bosh.pem"
echo "${BOSH_GW_PRIVATE_KEY}" > "${ssh_key}"
chmod 600 "${ssh_key}"
echo "Starting SSH tunnel, you may be prompted for your OS password..."
sudo true # prompt for password
sshuttle -e "ssh -i ${ssh_key} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" -r "${BOSH_GW_USER}@${BOSH_GW_HOST}" ${SSH_DESTINATION_CIDR} &
tunnel_pid="$!"

cleanup() {
  kill "${tunnel_pid}"
  rm -rf "${tmpdir}"
}
trap 'cleanup' EXIT

export BBR_BUILD_PATH=$(which bbr)

echo "Running DRATs..."
go get github.com/onsi/ginkgo/ginkgo
dep ensure
ginkgo -v --trace acceptance

echo "Successfully ran DRATs!"
