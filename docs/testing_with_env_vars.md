## Running DRATs using environment variables

In order to run DRATs using environment variables, you can use [`scripts/run_acceptance_tests_local.sh`](../scripts/testing_with_env_vars.md).

### Vanilla DRATs
For a vanilla run of DRATs, you need to set the following environment variables to run all of the default component test cases:
<table style="width:100%">
  <tr>
    <th>Environment Variable</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`CF_DEPLOYMENT_NAME`</td>
        <td>Name of the Cloud Foundry deployment to backup and restore</td>
    <tr>
    <tr>
        <td>`CF_API_URL`</td>
        <td>Cloud Foundry API URL</td>
    </tr>
    <tr>
        <td>`CF_ADMIN_USERNAME`</td>
        <td>Cloud Foundry API admin username</td>
    </tr>
    <tr>
        <td>`CF_ADMIN_PASSWORD`</td>
        <td>Cloud Foundry API admin password</td>
    </tr>
    <tr>
        <td>`CF_CREDHUB_CLIENT`</td>
        <td>Credhub for Cloud Foundry admin user</td>
    </tr>
    <tr>
        <td>`CF_CREDHUB_SECRET`</td>
        <td>Credhub for Cloud Foundry admin password</td>
    </tr>
    <tr>
        <td>`BOSH_ENVIRONMENT`</td>
        <td>URL of BOSH Director that deployed `CF_DEPLOYMENT_NAME`</td>
    </tr>
    <tr>
        <td>`BOSH_CLIENT`</td>
        <td>BOSH Director username</td>
    </tr>
    <tr>
        <td>`BOSH_CLIENT_SECRET`</td>
        <td>BOSH Director password</td>
    </tr>
    <tr>
        <td>`BOSH_CA_CERT`</td>
        <td>BOSH Director's CA cert content</td>
    </tr>
    <tr>
        <td>`BOSH_GW_HOST`</td>
        <td>Gateway host to use for BOSH SSH connection</td>
    </tr>
    <tr>
        <td>`BOSH_GW_USER`</td>
        <td>Gateway user to use for BOSH SSH connection</td>
    </tr>
    <tr>
        <td>`JUMPBOX_PRIVATE_KEY`</td>
        <td>Private key to use for BOSH SSH connection</td>
    </tr>
    <tr>
        <td>`BBR_BUILD_PATH`</td>
        <td>Path to BBR binary</td>
    </tr>
</table>

---

For further configure the run of DRATs, you can set the following environment variables:
<table style="width:100%">
  <tr>
    <th>Optional Environment Variable</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`DEFAULT_TIMEOUT_MINS`</td>
        <td>Timeout for commands run in the test (defaults to 15 minutes)</td>
    </tr>
    <tr>
        <td>`SSH_DESTINATION_CIDR`</td>
        <td>Default to "10.0.0.0/8"; change if your CF deployment is deployed in a different internal network range</td>
    <tr>
    <tr>
        <td>`DELETE_AND_REDEPLOY_CF`</td>
        <td>Set to "true" to have the CF deployment destroyed and redeployed from scratch during the test cycle. 
        **<span style="color:red">Exercise extreme care when using this option!</span>**</td>
    <tr>
    <tr>
        <td>`FOCUSED_SUITE_NAME`</td>
        <td>A regex matching the name(s) of test suites that you **do** want DRATs to run**</td>
    <tr>
    <tr>
        <td>`SKIP_SUITE_NAME`</td>
        <td>A regex matching the name(s) of test suites that you **do not** want DRATs to run**</td>
    <tr>
</table>

**If these variables are not set, all test suites will be run.

---

### DRATs with NFS Broker
To also run the test case for the optional component NFS Broker, set the following variables:
<table style="width:100%">
  <tr>
    <th>NFS Broker Environment Variable</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`INCLUDE_NFS_BROKER_TESTCASE`</td>
        <td>Set to "true" to run the NFS test case</td>
    <tr>
    <tr>
        <td>`NFS_SERVICE_NAME`</td>
        <td>Required to run the NFS test case</td>
    <tr>
    <tr>
        <td>`NFS_PLAN_NAME`</td>
        <td>Required to run the NFS test case</td>
    <tr>
    <tr>
        <td>`NFS_CREATE_SERVICE_BROKER`</td>
        <td>Set to "true" to register the NFS service broker in the NFS test case</td>
    <tr>
    <tr>
        <td>`NFS_BROKER_USER`</td>
        <td>Required to register the NFS service broker when running the NFS test case</td>
    <tr>
    <tr>
        <td>`NFS_BROKER_PASSWORD`</td>
        <td>Required to register the NFS service broker when running the NFS test case</td>
    <tr>
    <tr>
        <td>`NFS_BROKER_URL`</td>
        <td>Required to register the NFS service broker when running the NFS test case</td>
    <tr>
</table>

---

### DRATs with SMB Broker
To also run the test case for the optional component SMB Broker, set the following variables:
<table style="width:100%">
  <tr>
    <th>SMB Broker Environment Variable</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`INCLUDE_SMB_BROKER_TESTCASE`</td>
        <td>Set to "true" to run the SMB test case</td>
    <tr>
    <tr>
        <td>`SMB_SERVICE_NAME`</td>
        <td>Required to run the SMB test case</td>
    <tr>
    <tr>
        <td>`SMB_PLAN_NAME`</td>
        <td>Required to run the SMB test case</td>
    <tr>
    <tr>
        <td>`SMB_CREATE_SERVICE_BROKER`</td>
        <td>Set to "true" to register the SMB service broker in the SMB test case</td>
    <tr>
    <tr>
        <td>`SMB_BROKER_USER`</td>
        <td>Required to register the SMB service broker when running the SMB test case</td>
    <tr>
    <tr>
        <td>`SMB_BROKER_PASSWORD`</td>
        <td>Required to register the SMB service broker when running the SMB test case</td>
    <tr>
    <tr>
        <td>`SMB_BROKER_URL`</td>
        <td>Required to register the SMB service broker when running the SMB test case</td>
    <tr>
</table>
