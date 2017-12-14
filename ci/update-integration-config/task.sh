#!/usr/bin/env bash

set -eu

export CF_DEPLOYMENT_NAME
export INTEGRATION_CONFIG_FILE_PATH
export VARS_STORE_FILE_PATH
export BBL_STATE_DIR_PATH
export SYSTEM_DOMAIN

cf_deployment_name="${CF_DEPLOYMENT_NAME}"
cf_api_url="https://api.${SYSTEM_DOMAIN}"
cf_admin_username=admin
cf_admin_password=`bosh interpolate --path=/cf_admin_password vars-store/${VARS_STORE_FILE_PATH}`
bosh_environment=`bosh interpolate --path=/external_ip <(bbl --state-dir=bbl-state-store/${BBL_STATE_DIR_PATH} bosh-deployment-vars)`
bosh_client=`bbl --state-dir=bbl-state-store/${BBL_STATE_DIR_PATH} director-username`
bosh_client_secret=`bbl --state-dir=bbl-state-store/${BBL_STATE_DIR_PATH} director-password`
bosh_ca_cert=`bbl --state-dir=bbl-state-store/${BBL_STATE_DIR_PATH} director-ca-cert`
ssh_proxy_user="jumpbox"
ssh_proxy_host=${bosh_environment}
ssh_proxy_cidr="10.0.0.0/8"
ssh_proxy_private_key=`bbl --state-dir=bbl-state-store/${BBL_STATE_DIR_PATH} ssh-key`
nfs_broker_password=`bosh interpolate --path=/nfs-broker-password vars-store/${VARS_STORE_FILE_PATH} || echo ""`
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
