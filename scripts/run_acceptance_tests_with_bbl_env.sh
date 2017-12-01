#!/bin/bash
: "${CF_VARS_STORE_PATH:="cf-deployment-variables.yml"}"
: "${BOSH_CLI_NAME:="bosh"}"

pushd $1
    export CF_ADMIN_PASSWORD=$(${BOSH_CLI_NAME} interpolate --path /cf_admin_password ${CF_VARS_STORE_PATH})
    eval "$(bbl print-env | sed '$ d'| sed '$ d'| sed '$ d')"
    export BOSH_GW_USER="jumpbox"
    export BOSH_GW_HOST=$(bbl director-address | cut -d'/' -f3 | cut -d':' -f1)
    export BOSH_GW_PRIVATE_KEY_CONTENTS="$(bbl director-ssh-key)"

    CF_DOMAIN=$(jq .lb.domain bbl-state.json -r)
    export CF_API_URL="https://api.${CF_DOMAIN}"
    if grep "nfs-broker-password" ${CF_VARS_STORE_PATH}>/dev/null; then
        export NFS_SERVICE_NAME="nfs"
        export NFS_PLAN_NAME="Existing"
        export NFS_BROKER_USER="nfs-broker"
        export NFS_BROKER_PASSWORD=$(${BOSH_CLI_NAME} interpolate --path /nfs-broker-password ${CF_VARS_STORE_PATH})
        export NFS_BROKER_URL="http://nfs-broker.${CF_DOMAIN}"
    else
        echo "Skipping cf-nfsrboker testcase because nfs-broker-password is not present in ${CF_VARS_STORE_PATH}"
        export SKIP_SUITE_NAME="cf-nfsbroker"
    fi
popd

echo "Running DRATs locally"
. ./scripts/run_acceptance_tests_local.sh
