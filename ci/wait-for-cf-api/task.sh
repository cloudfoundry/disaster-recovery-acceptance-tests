#!/usr/bin/env bash

call_cf_api() {
    cf api "${CF_API_URL}" --skip-ssl-validation
}

call_uaa_api() {
    echo "curl -k ${CF_UAA_URL}"
    curl -k "${CF_UAA_URL}"
}

set -e
curl "${CF_API_URL}" --fail --retry "${RETRY_COUNT}" --insecure
set +e

for i in $(seq 1 "${RETRY_COUNT}"); do
    call_cf_api
    exit_code=$?

    if [[ ${exit_code} -eq 0 ]]; then
      call_uaa_api
      exit_code=$?
      echo "exit code: $exit_code"
      if [[ ${exit_code} -eq 0 ]]; then
          break
      fi
    fi

    sleep 15
done

if [[ ${exit_code} -ne 0 ]]; then
   echo "Failed to successfully connect to CF API after ${RETRY_COUNT} tries"
   exit 1
fi

set -e

# i is unused
# shellcheck disable=SC2034
for i in $(seq 1 "${RETRY_COUNT}"); do
    sleep 15

    call_cf_api
    call_uaa_api
done
