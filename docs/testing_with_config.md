## Running DRATS with Integration Config

In order to run DRATs using an integration config file, you can use [`scripts/run_acceptance_tests_local_with_config.sh`](../scripts/run_acceptance_tests_local_with_config.sh) or on the jumpbox using `ginkgo`.


### From a jumpbox
1. Create an integration config json file e.g. [integration_config.json](../ci/integration_config.json):
1. Export the following environment variables:
   * `BBR_BUILD_PATH`: Path to your `bbr` binary
   * `CONFIG`: Path to your `integration_config.json`
1. Run the tests
   ```bash
   $> dep ensure
   $> ginkgo -v --trace acceptance
   ```

### Locally
1. Create an integration config json file e.g. [integration_config.json](../ci/integration_config.json):
1. Setup the following environment variables:
   * `BBR_BUILD_PATH`
1. Add the following additional properties to the integration config:
   *  `ssh_proxy_user`: The jumpbox user
   *  `ssh_proxy_host`: The jumpbox host
   *  `ssh_proxy_private_key`: The jumpbox private key
1. Run [`scripts/run_acceptance_tests_local_with_config.sh`](../scripts/run_acceptance_tests_local_with_config.sh)

### Required Integration Config Properties
To run all of the default test cases in DRATs, you must set the following integration config properties:

<table style="width:100%">
  <tr>
    <th>Required Integration Config Properties</th>
    <th>Usage</th>
  </tr>
  <tr>
    <td>`cf_deployment_name`</td>
    <td>Name of the Cloud Foundry deployment to backup and restore</td>
  </tr>
  <tr>
    <td>`cf_api_url`</td>
    <td>Cloud Foundry API URL</td>
  </tr>
  <tr>
    <td>`cf_admin_username`</td>
    <td>Cloud Foundry API admin user</td>
  </tr>
  <tr>
    <td>`cf_admin_password`</td>
    <td>Cloud Foundry API admin password</td>
  </tr>
    <tr>
      <td>`credhub_client_name`</td>
      <td>Credhub for Cloud Foundry admin user</td>
    </tr>
    <tr>
      <td>`credhub_client_secret`</td>
      <td>Credhub for Cloud Foundry admin password</td>
    </tr>
  <tr>
    <td>`bosh_environment`</td>
    <td>URL of BOSH Director which has deployed the above Cloud Foundries</td>
  </tr>
  <tr>
    <td>`bosh_client`</td>
    <td>BOSH Director username</td>
  </tr>
  <tr>
    <td>`bosh_client`</td>
    <td>BOSH Director password</td>
  </tr>
  <tr>
    <td>`bosh_ca_cert`</td>
    <td>BOSH Director's CA cert content</td>
  </tr>
  <tr>
    <td>`include_TESTCASE_NAME`</td>
    <td>Flag for whether to run a given testcase; if omitted, it defaults to false</td>
  </tr>
</table>

---

### Optional Integration Config Properties
For further configure the run of DRATs, you can set the following integration config properties:
<table style="width:100%">
  <tr>
    <th>Optional Integration Config Property</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`timeout_in_minutes`</td>
        <td>Timeout for commands run in the test (defaults to 15 minutes)</td>
    </tr>
    <tr>
        <td>`ssh_destination_cidr`</td>
        <td>Default to "10.0.0.0/8"; change if your CF deployment is deployed in a different internal network range</td>
    <tr>
    <tr>
        <td>`delete_and_redeploy_cf`</td>
        <td>Set to "true" to have the CF deployment destroyed and redeployed from scratch during the test cycle. 
        **<span style="color:red">Exercise extreme care when using this option!</span>**</td>
    <tr>
</table>

**If these variables are not set, all test suites will be run.

---

### DRATs with NFS Broker
To also run the test case for the optional component NFS Broker, set the following properties:
<table style="width:100%">
  <tr>
    <th>NFS Broker Integration Config Property</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`nfs_service_name`</td>
        <td>Required to run the NFS test case</td>
    <tr>
    <tr>
        <td>`nfs_plan_name`</td>
        <td>Required to run the NFS test case</td>
    <tr>
    <tr>
        <td>`nfs_create_service_broker`</td>
        <td>Set to "true" to register the NFS service broker in the NFS test case</td>
    <tr>
    <tr>
        <td>`nfs_broker_user`</td>
        <td>Required to register the NFS service broker when running the NFS test case</td>
    <tr>
    <tr>
        <td>`nfs_broker_password`</td>
        <td>Required to register the NFS service broker when running the NFS test case</td>
    <tr>
    <tr>
        <td>`nfs_broker_url`</td>
        <td>Required to register the NFS service broker when running the NFS test case</td>
    <tr>
</table>

---

### DRATs with SMB Broker
To also run the test case for the optional component SMB Broker, set the following properties:
<table style="width:100%">
  <tr>
    <th>SMB Broker Integration Config Property</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`smb_service_name`</td>
        <td>Required to run the SMB test case</td>
    <tr>
    <tr>
        <td>`smb_plan_name`</td>
        <td>Required to run the SMB test case</td>
    <tr>
    <tr>
        <td>`smb_create_service_broker`</td>
        <td>Set to "true" to register the SMB service broker in the SMB test case</td>
    <tr>
    <tr>
        <td>`smb_broker_user`</td>
        <td>Required to register the SMB service broker when running the SMB test case</td>
    <tr>
    <tr>
        <td>`smb_broker_password`</td>
        <td>Required to register the SMB service broker when running the SMB test case</td>
    <tr>
    <tr>
        <td>`smb_broker_url`</td>
        <td>Required to register the SMB service broker when running the SMB test case</td>
    <tr>
</table>


### DRATs with Selective Backup
If you have deployed your cf-deployment to selectively backup blobs, set the following variables:
<table style="width:100%">
  <tr>
    <th>Selective Backup Environment Variable</th>
    <th>Usage</th>
  </tr>
    <tr>
        <td>`selective_backup`</td>
        <td>Set to "true" to run the EnsureAfterSelectiveRestore test case step</td>
    <tr>
    <tr>
        <td>`selective_backup_type`</td>
        <td>set to droplets or droplets_and_packages</td>
    <tr>
</table>
