#!/bin/bash
: "${CF_VARS_STORE_PATH:="cf-deployment-variables.yml"}"
: "${BOSH_CLI_NAME:="bosh"}"

pushd $1
    export CF_ADMIN_PASSWORD=$(${BOSH_CLI_NAME} interpolate --path /cf_admin_password ${CF_VARS_STORE_PATH})
    export BOSH_CLIENT_SECRET=$(bbl director-password)
    export BOSH_CA_CERT="$(bbl director-ca-cert)"
    export BOSH_ENVIRONMENT=$(${BOSH_CLI_NAME} interpolate --path /external_ip <(bbl bosh-deployment-vars))
    export BOSH_GW_USER="jumpbox"
    export BOSH_GW_HOST=$(${BOSH_CLI_NAME} interpolate --path /external_ip <(bbl bosh-deployment-vars))
    export BOSH_GW_PRIVATE_KEY_CONTENTS="$(bbl ssh-key)"
    export BOSH_CLIENT="admin"
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
        if [ -z ${SKIP_SUITE_NAME} ]; then
            export SKIP_SUITE_NAME="cf-nfsbroker"
        else
            export SKIP_SUITE_NAME="(${SKIP_SUITE_NAME})|cf-nfsbroker"
        fi
    fi

    if grep "azurefile-broker-password" ${CF_VARS_STORE_PATH}>/dev/null; then
        export SMB_SERVICE_NAME="smb"
        export SMB_PLAN_NAME="Existing"
        export SMB_BROKER_USER="admin"
        export SMB_BROKER_PASSWORD=$(${BOSH_CLI_NAME} interpolate --path /azurefile-broker-password ${CF_VARS_STORE_PATH})
        export SMB_BROKER_URL="http://azurefile-broker.${CF_DOMAIN}"
    else
        echo "Skipping cf-smbbroker testcase because azurefile-broker-password is not present in ${CF_VARS_STORE_PATH}"
        if [ -z ${SKIP_SUITE_NAME} ]; then
            export SKIP_SUITE_NAME="cf-smbbroker"
        else
            export SKIP_SUITE_NAME="(${SKIP_SUITE_NAME})|cf-smbbroker"
        fi

    fi
popd

echo "Running DRATs locally"
. ./scripts/run_acceptance_tests_local.sh
