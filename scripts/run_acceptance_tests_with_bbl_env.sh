#!/bin/bash

# INCLUDE_NFS_BROKER_TESTCASE="true"
# NFS_CREATE_SERVICE_BROKER="true"
# INCLUDE_SMB_BROKER_TESTCASE="true"
# SMB_CREATE_SERVICE_BROKER="true"

bbl_directory=$1

function find_cred {
  credhub find -j -n "$1" | jq .credentials[0].name | xargs credhub get -j -n | jq -r .value
}

pushd "$bbl_directory" || exit
  eval "$(bbl print-env)"
  export CF_ADMIN_PASSWORD=$(find_cred cf_admin_password)
  CF_DOMAIN=$(jq .lb.domain bbl-state.json -r)
  export CF_API_URL="https://api.${CF_DOMAIN}"

  if [[ "${INCLUDE_NFS_BROKER_TESTCASE}" = "true" ]]; then
    export NFS_SERVICE_NAME="nfs"
    export NFS_PLAN_NAME="Existing"

    if [[ "${NFS_CREATE_SERVICE_BROKER}" = "true" ]]; then
      export NFS_BROKER_USER="nfs-broker"
      export NFS_BROKER_PASSWORD=$(find_cred nfs-broker-password)
      export NFS_BROKER_URL="http://nfs-broker.${CF_DOMAIN}"
    fi
  else
      echo "Skipping cf-nfsbroker testcase because INCLUDE_NFS_BROKER_TESTCASE is not set to true"
  fi

  if [[ "${INCLUDE_SMB_BROKER_TESTCASE}" = "true" ]]; then
    export SMB_SERVICE_NAME="smb"
    export SMB_PLAN_NAME="Existing"

    if [[ "${SMB_CREATE_SERVICE_BROKER}" = "true" ]]; then
      export SMB_BROKER_USER="admin"
      export SMB_BROKER_PASSWORD=$(find_cred smb-broker-password)
      export SMB_BROKER_URL="http://smbbroker.${CF_DOMAIN}"
    fi
  else
      echo "Skipping cf-smbbroker testcase because INCLUDE_SMB_BROKER_TESTCASE is not set to true"
  fi
popd || exit

echo "Running DRATs locally"
. ./scripts/run_acceptance_tests_local.sh

