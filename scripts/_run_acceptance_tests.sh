#!/usr/bin/env bash

set -eu

export CF_DEPLOYMENT_NAME
export CF_API_URL
export CF_ADMIN_USERNAME
export CF_ADMIN_PASSWORD
export NFS_SERVICE_NAME
export NFS_PLAN_NAME
export NFS_BROKER_USER
export NFS_BROKER_PASSWORD
export NFS_BROKER_URL
export BOSH_CA_CERT
export BOSH_CLIENT
export BOSH_CLIENT_SECRET
export BOSH_ENVIRONMENT
export BBR_BUILD_PATH
export FOCUSED_SUITE_NAME
export SKIP_SUITE_NAME
export DELETE_AND_REDEPLOY_CF

go get github.com/onsi/ginkgo/ginkgo
dep ensure
ginkgo -v --trace acceptance
