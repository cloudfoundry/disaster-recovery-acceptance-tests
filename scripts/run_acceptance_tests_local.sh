#!/bin/bash

set -eu -o pipefail

# ENV
: "${BOSH_CLIENT:?}"
: "${BOSH_CLIENT_SECRET:?}"
: "${BOSH_CA_CERT:?}"
: "${BOSH_GW_HOST:?}"
: "${BOSH_GW_USER:?}"
: "${BOSH_GW_PRIVATE_KEY_CONTENTS:?}"
: "${CF_ADMIN_PASSWORD:?}"
: "${CF_API_URL:?}"
: "${GOPATH:?}"
: "${CF_DEPLOYMENT_NAME:="cf"}"
: "${CF_ADMIN_USERNAME:="admin"}"
: "${BOSH_ENVIRONMENT:?}"
: "${SSH_DESTINATION_CIDR:="10.0.0.0/8"}"
# The following params are optional
: "${NFS_SERVICE_NAME:=""}"
: "${NFS_PLAN_NAME:=""}"
: "${NFS_BROKER_USER:=""}"
: "${NFS_BROKER_PASSWORD:=""}"
: "${NFS_BROKER_URL:=""}"
: "${SMB_SERVICE_NAME:=""}"
: "${SMB_PLAN_NAME:=""}"
: "${SMB_BROKER_USER:=""}"
: "${SMB_BROKER_PASSWORD:=""}"
: "${SMB_BROKER_URL:=""}"
: "${SKIP_SUITE_NAME:=""}"

tmpdir="$( mktemp -d /tmp/run-drats.XXXXXXXXXX )"

ssh_key="${tmpdir}/bosh.pem"
echo "${BOSH_GW_PRIVATE_KEY_CONTENTS}" > "${ssh_key}"
chmod 600 "${ssh_key}"
echo "Starting SSH tunnel, you may be prompted for your OS password..."
sudo true # prompt for password
# ssh -i "${ssh_key}" -o 'StrictHostKeyChecking no' -o 'UserKnownHostsFile /dev/null' ${BOSH_GW_USER}@${BOSH_GW_HOST}
sshuttle -e "ssh -i ${ssh_key} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" -r "${BOSH_GW_USER}@${BOSH_GW_HOST}" ${SSH_DESTINATION_CIDR} &
tunnel_pid="$!"

cleanup() {
  kill "${tunnel_pid}"
  rm -rf "${tmpdir}"
}
trap 'cleanup' EXIT

if [ -n "${BOSH_CA_CERT}" ]; then
  export BOSH_CERT_PATH="${tmpdir}/bosh.ca"
  echo "${BOSH_CA_CERT}" > "${BOSH_CERT_PATH}"
fi

export BBR_BUILD_PATH="${BBR_BUILD_PATH:-$(which bbr)}"
export BOSH_URL="${BOSH_ENVIRONMENT}"

echo "Running DRATs..."
. ./scripts/_run_acceptance_tests.sh

echo "Successfully ran DRATs!"
