# disaster-recovery-acceptance-tests (DRATS)

Tests if Cloud Foundry can be backed up and restored. The tests will back up from and restore to `CF_DEPLOYMENT_NAME`.

## Prerequisites
1. Install golang, per golang.org.   
1. Install `ginkgo`
   ```bash
   go get github.com/onsi/ginkgo/ginkgo
   ```
1. Install `dep`

## Running DRATS with Envrironment Variables

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
1. If you wish to run DRATS against a director deployed with `bbl`, run `scripts/run_acceptance_tests_with_bbl_env.sh <path-to-bbl-state-dir>`.
    * Set `CF_VARS_STORE_PATH` to the path to the CF vars-store file.
    * Set `BOSH_CLI_NAME` to the name of the BOSH CLI executable on your machine if it isn't `bosh`.
    * The script will search for NFS broker credentials in the CF vars-store file and will skip the NFS test case if those credentials are not present.

## Running DRATS with Integration Config

#### From a jumpbox
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
     "bosh_ca_cert": "-----BEGIN CERTIFICATE------\n...\n------END CERTIFICATE-----"
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
#### Locally
1. Follow first two steps in jumpbox instructions.
1. Setup the following environment variables:
   *  `SSH_DESTINATION_CIDR`
   *  `BOSH_GW_HOST`
   *  `BOSH_GW_USER`
   *  `BOSH_GW_PRIVATE_KEY_CONTENTS`
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

#### Optional Variables
* `nfs_service_name` - Environment variable required to run NFS test case
* `nfs_plan_name` - Environment variable required to run NFS test case
* `nfs_broker_user` - Environment variable required to run NFS test case
* `nfs_broker_password` - Environment variable required to run NFS test case
* `nfs_broker_url` - Environment variable required to run NFS test case

### Focusing/Skipping a test suite

Run DRATS as usual but set the environment variable `FOCUSED_SUITE_NAME` and/or `SKIP_SUITE_NAME` to a regex matching the name(s) of test suites. Only those suites that either match `FOCUSED_SUITE_NAME` or don't match `SKIP_SUITE_NAME` will be run.  Leaving either of these unset is supported.

If these variables are not set, all test suites returned by [`testcases.OpenSourceTestCases()`](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/blob/master/testcases/testcase_helper.go#L9) will be run.

## Test Structure

The system tests do the following:

1. Calls `BeforeBackup(common.Config)` on all provided TestCases (to e.g. push unique apps to the environment to be backed up).
1. Backs up the `CF_DEPLOYMENT_NAME` Cloud Foundry deployment.
1. Calls `AfterBackup(common.Config)` on all provided TestCases.
1. Restores to the `CF_DEPLOYMENT_NAME` Cloud Foundry deployment.
1. Calls `AfterRestore(common.Config)` on all provided TestCases (to e.g. check the apps pushed are present in the restored environment).
1. Calls `Cleanup(common.Config)` on all provided TestCases (to e.g. clean up the apps from the backup environment).

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

We have shared a [task](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/tree/master/ci/drats) to run DRATS with your CI. The task establishes an SSH tunnel using [`sshuttle`](http://sshuttle.readthedocs.io) so that it can run from outside the network. Note that this task needs a privileged container.
