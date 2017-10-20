#!/bin/bash

pushd $1
    export CF_ADMIN_PASSWORD=$(bosh-cli interpolate --path /cf_admin_password cf-deployment-variables.yml)
    export BOSH_CLIENT_SECRET=$(bbl director-password)
    export BOSH_CA_CERT="$(bbl director-ca-cert)"
    export BOSH_ENVIRONMENT=$(bosh-cli interpolate --path /external_ip <(bbl bosh-deployment-vars))
    export BOSH_GW_USER="jumpbox"
    export BOSH_GW_HOST=$(bosh-cli interpolate --path /external_ip <(bbl bosh-deployment-vars))
    export BOSH_GW_PRIVATE_KEY_CONTENTS="$(bbl ssh-key)"
    export BOSH_CLIENT="admin"
    CF_DOMAIN=$(jq .lb.domain bbl-state.json -r)
    export CF_API_URL="api.${CF_DOMAIN}"
    export NFS_SERVICE_NAME="nfs"
    export NFS_PLAN_NAME="Existing"
    export NFS_BROKER_USER="nfs-broker"
    export NFS_BROKER_PASSWORD=$(bosh-cli interpolate --path /nfs-broker-password cf-deployment-variables.yml)
    export NFS_BROKER_URL="http://nfs-broker.${CF_DOMAIN}"
popd

echo "Running DRATs locally"
. ./scripts/run_acceptance_tests_local.sh

