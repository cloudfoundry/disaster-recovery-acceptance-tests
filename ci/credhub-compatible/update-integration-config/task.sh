#!/usr/bin/env bash

set -euo pipefail

get_password_from_credhub() {
  local bosh_manifest_password_variable_name=$1
  echo $(credhub find -j -n ${bosh_manifest_password_variable_name} | jq -r .credentials[].name | xargs credhub get -j -n | jq -r .value)
}

setup_bosh_env_vars() {
  pushd "bbl-state/${BBL_STATE_DIR}"
    eval "$(bbl print-env)"
  popd
}

setup_bosh_env_vars

cf_deployment_name="${CF_DEPLOYMENT_NAME}"
cf_api_url="https://api.${SYSTEM_DOMAIN}"
cf_admin_username=admin
cf_admin_password=$(get_password_from_credhub cf_admin_password)
bosh_environment="$BOSH_ENVIRONMENT"
bosh_client="$BOSH_CLIENT"
bosh_client_secret="$BOSH_CLIENT_SECRET"
bosh_ca_cert="$BOSH_CA_CERT"
ssh_proxy_user="jumpbox"
ssh_proxy_host=$(bbl --state-dir "bbl-state/$BBL_STATE_DIR" jumpbox-address)
ssh_proxy_cidr="10.0.0.0/8"
ssh_proxy_private_key="$(cat "$JUMPBOX_PRIVATE_KEY")"
nfs_broker_password=$(get_password_from_credhub nfs-broker-password || echo "")
nfs_service_name="nfs"
nfs_plan_name="Existing"
nfs_broker_user="nfs-broker"
nfs_broker_url="http://nfs-broker.${SYSTEM_DOMAIN}"

configs=( cf_deployment_name
        cf_api_url
        cf_admin_username
        cf_admin_password
        bosh_environment
        bosh_client
        bosh_client_secret
        bosh_ca_cert
        ssh_proxy_user
        ssh_proxy_host
        ssh_proxy_cidr
        ssh_proxy_private_key
        nfs_broker_password
        nfs_service_name
        nfs_plan_name
        nfs_broker_user
        nfs_broker_url )

integration_config=`cat integration-configs/${INTEGRATION_CONFIG_FILE_PATH}`

for config in "${configs[@]}"
do
  integration_config=$(echo ${integration_config} | jq ".${config}=\"${!config}\"")
done

if [ -z ${nfs_broker_password} ]; then
  integration_config=$(echo ${integration_config} | jq '."include_cf-nfsbroker"=false')
fi

echo "${integration_config}" > integration-configs/${INTEGRATION_CONFIG_FILE_PATH}

cp -Tr integration-configs updated-integration-configs
