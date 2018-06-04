# disaster-recovery-acceptance-tests (DRATS)

Tests if Cloud Foundry can be backed up and restored. The tests will back up from and restore to `CF_DEPLOYMENT_NAME`.

## Prerequisites
1. Install golang, per golang.org.   
1. Install `ginkgo`
   ```bash
   go get github.com/onsi/ginkgo/ginkgo
   ```
1. Install `dep`

## Running DRATS with Environment Variables

1. Spin up a Cloud Foundry deployment.
    * CF on BOSH Lite is supported.
    * [cf-deployment](https://github.com/cloudfoundry/cf-deployment) is supported. Ensure you apply the [backup-restore opsfile](https://github.com/cloudfoundry/cf-deployment/blob/master/operations/experimental/enable-backup-restore.yml) at deploy time to ensure the backup and restore scripts are enabled. This will also deploy a backup restore VM, on which the [Backup and Restore SDK](https://github.com/cloudfoundry-incubator/backup-and-restore-sdk-release) is deployed.
1. Run `scripts/run_acceptance_tests_local.sh` with the following environment variables set:
    * `CF_DEPLOYMENT_NAME` - name of the Cloud Foundry deployment to backup and restore
    * `CF_API_URL` - Cloud Foundry api url
    * `CF_ADMIN_USERNAME` - Cloud Foundry api admin user
    * `CF_ADMIN_PASSWORD` - Cloud Foundry api admin password
    * `BOSH_ENVIRONMENT` - URL of BOSH Director which has deployed the above Cloud Foundries
    * `BOSH_CLIENT` - BOSH Director username
    * `BOSH_CLIENT_SECRET` - BOSH Director password
    * `BOSH_CA_CERT` - BOSH Director's CA cert content
    * `BOSH_GW_HOST` - Gateway host to use for BOSH SSH connection
    * `BOSH_GW_USER` - Gateway user to use for BOSH SSH connection
    * `BOSH_GW_PRIVATE_KEY_CONTENTS` - Private key to use for BOSH SSH connection
    * `BBR_BUILD_PATH` - path to BBR binary
    * `DEFAULT_TIMEOUT_MINS` - timeout for commands run in the test. Defaults to 15 minutes.
1. The following environment variables are optional and could be set depending on test configuration:
    * `SSH_DESTINATION_CIDR` - Default to "10.0.0.0/8"; change if your cf-deployment is deployed in a different internal network range
    * `NFS_SERVICE_NAME` - Environment variable required to run NFS test case
    * `NFS_PLAN_NAME` - Environment variable required to run NFS test case
    * `NFS_BROKER_USER` - Environment variable required to run NFS test case
    * `NFS_BROKER_PASSWORD` - Environment variable required to run NFS test case
    * `NFS_BROKER_URL` - Environment variable required to run NFS test case
    * `SMB_SERVICE_NAME` - Environment variable required to run SMB test case
    * `SMB_PLAN_NAME` - Environment variable required to run SMB test case
    * `SMB_BROKER_USER` - Environment variable required to run SMB test case
    * `SMB_BROKER_PASSWORD` - Environment variable required to run SMB test case
    * `SMB_BROKER_URL` - Environment variable required to run SMB test case
    * `DELETE_AND_REDEPLOY_CF` - set to "true" to have the CF deployment destroyed and redeployed from scratch during the test cycle. **<span style="color:red"> Exercise extreme care when using this option!</span>**
1. If you wish to run DRATS against a director deployed with `bbl`, run `scripts/run_acceptance_tests_with_bbl_env.sh <path-to-bbl-state-dir>`.
    * Set `CF_VARS_STORE_PATH` to the path to the CF vars-store file.
    * Set `BOSH_CLI_NAME` to the name of the BOSH CLI executable on your machine if it isn't `bosh`.
    * The script will search for NFS/SMB broker credentials in the CF vars-store file and will skip the NFS/SMB test case if those credentials are not present.

### Focusing/Skipping a test suite

Run DRATS as usual but set the environment variable `FOCUSED_SUITE_NAME` and/or `SKIP_SUITE_NAME` to a regex matching the name(s) of test suites. Only those suites that both match `FOCUSED_SUITE_NAME` and don't match `SKIP_SUITE_NAME` will be run. Leaving either of these unset is supported.

If these variables are not set, all test suites returned by [`testcases.OpenSourceTestCases()`](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/blob/master/testcases/testcase_helper.go#L9) will be run.

## Running DRATS with Integration Config

### From a jumpbox
1. Create an integration config json file, for example:
   ```bash
   cat > integration_config.json <<EOF
    {
     "cf_api_url": "api.<cf_system_domain>",
     "cf_deployment_name": "cf",
     "cf_admin_username": "admin",
     "cf_admin_password": "<cf_admin_password>",
     "bosh_environment": "https://<bosh_director_ip>:25555",
     "bosh_client": "admin",
     "bosh_client_secret": "<bosh_admin_password>",
     "bosh_ca_cert": "-----BEGIN CERTIFICATE------\n...\n------END CERTIFICATE-----",
     "include_cf-app": true,
     "include_cf-networking": true,
     "include_cf-uaa": true,
     "include_cf-credhub": true,
     "include_cf-routing": true,
     "include_app-uptime": true
   }
   EOF
   export CONFIG=$PWD/integration_config.json
   ```
1. Setup the following environment variables:
   * `BBR_BUILD_PATH`
1. [Optional] Change the default timeout by setting `DEFAULT_TIMEOUT_MINS`.
1. Run the tests
   ```bash
   dep ensure
   ginkgo -v --trace acceptance
   ```
### Locally
1. Follow first two steps in jumpbox instructions.
1. Add the following additional properties to the integration config:
   *  `ssh_proxy_user`
   *  `ssh_proxy_host`
   *  `ssh_proxy_private_key`
   *  `ssh_proxy_cidr`
1. Run `scripts/run_acceptance_tests_local_with_config.sh`

### Integration Config Variables
* `cf_deployment_name` - name of the Cloud Foundry deployment to backup and restore
* `cf_api_url` - Cloud Foundry api url
* `cf_admin_username` - Cloud Foundry api admin user
* `cf_admin_password` - Cloud Foundry api admin password
* `bosh_environment` - URL of BOSH Director which has deployed the above Cloud Foundries
* `bosh_client` - BOSH Director username
* `bosh_client_secret` - BOSH Director password
* `bosh_ca_cert` - BOSH Director's CA cert content
* `include_<testcase-name>` - Flag for whether to run a given testcase. If omitted defaults to false

#### Optional Variables
* `nfs_service_name` - Environment variable required to run NFS test case
* `nfs_plan_name` - Environment variable required to run NFS test case
* `nfs_broker_user` - Environment variable required to run NFS test case
* `nfs_broker_password` - Environment variable required to run NFS test case
* `nfs_broker_url` - Environment variable required to run NFS test case
* `timeout_in_minutes` - Default ginkgo `Eventually` timeout. Defaults to 15
* `delete_and_redeploy_cf` - Destroy and redeploy the cf between after backup and restore. Defaults to false.

## Test Structure

The system tests do the following:

1. Sets up a temporary local working directory for storing the backup artifact, and CF_HOME directories for all the test cases.
1. Calls `BeforeBackup(common.Config)` on all provided TestCases (to e.g. push unique apps to the environment to be backed up).
1. Backs up the `CF_DEPLOYMENT_NAME` Cloud Foundry deployment.
1. Calls `AfterBackup(common.Config)` on all provided TestCases.
1. Restores to the `CF_DEPLOYMENT_NAME` Cloud Foundry deployment.
1. Calls `AfterRestore(common.Config)` on all provided TestCases (to e.g. check the apps pushed are present in the restored environment).
1. Calls `Cleanup(common.Config)` on all provided TestCases (to e.g. clean up the apps from the backup environment). It will do this even if an error or failure occurred in a previous step
1. Cleans up the temporary directories created in the setup
1. If an error occurred during a `bbr backup` command, DRATS runs `bbr backup-cleanup` to remove temporary bbr artifacts from your deployment (which would otherwise cause subsequent DRATS runs to fail)

## Extending DRATS

DRATS runs a collection of test cases against a Cloud Foundry deployment.

Test cases should be used for checking that CF components' data has been backed up and restored correctly – e.g. if your release backs up a table in a database, that the table can be altered and is then restored to its original state.

**Backup and restore of apps is covered by the existing CAPI test case.** No new test cases are needed for this – unless you're writing CAPI backup and restore scripts, app backup and restore can be assumed to work.

To add extra test cases, create a new TestCase that follows the [TestCase interface](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/blob/master/runner/testcase.go).

The methods that need to be implemented are `BeforeBackup(common.Config)`, `AfterBackup(common.Config)`, `AfterRestore(common.Config)` and `Cleanup(common.Config)`.

* `BeforeBackup(common.Config)` runs before the backup is taken, and should create state in the Cloud Foundry deployment to be backed up.
* `AfterBackup(common.Config)` runs after the backup is complete but before the restore is started. If were monitoring e.g. app uptime during the backup you could use this step to stop monitoring knowing that backup definitely finished.
* `AfterRestore(common.Config)` runs after the restore is complete, and should assert that the state in the restored Cloud Foundry deployment matches that created in `BeforeBackup(common.Config)`.
* `Cleanup(common.Config)` should clean up the state created in the Cloud Foundry deployment to be backed up.

`common.Config` contains the config for the BOSH Director and for the CF deployments to backup and restore.

1. Create a new TestCase in `test_cases`
1. In `testcases/testcase_helper.go`, initialise the TestCase and add it to the slice returned by `OpenSourceTestCases()`

## Running DRATs in your CI

We have shared a [task](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/tree/master/ci/drats) to run DRATS with your CI. This task reads credentials from environment variables. If you prefer a CATs like integration_config, we have another CI [task](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/tree/master/ci/drats-with-integration-config) for you.
Both run-drats tasks establish an SSH tunnel using [`sshuttle`](http://sshuttle.readthedocs.io) so that they can run from outside the network. Note the tasks will need to be run from a privileged container.
You can also find our pipeline definition here: [pipeline](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/tree/master/ci/pipeline.yml)

## Debugging your DRATS run

DRATS runs multiple interwoven test cases (for app uptime and each of the components under test) so it can be a little tricky to work out what's gone wrong when there's an error or failure. Here are some tips on investigating DRATS failures - please PR in additions to this doc if you think of more tips that might help other teams!

1. The `bbr backup-cleanup` command runs if the test run errored during the `bbr backup` step. If you see an error in the `backup-cleanup` step, it's likely that a similar problem happened in the `backup` step which caused the original failure - scroll up to see.
1. The easiest way to see where the failure / error happened is to look for the nearest `STEP` statement in the logs
