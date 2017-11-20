# disaster-recovery-acceptance-tests (DRATS)

Tests if Cloud Foundry can be backed up and restored. The tests will back up from and restore to `cf_deployment_name`.

## Running DRATS

### Prerequisites

You will need a Cloud Foundry deployment. This can be CF on BOSH Lite or 
[cf-deployment](https://github.com/cloudfoundry/cf-deployment). If you are using `cf-deployment`, 
ensure you apply the 
[backup-restore opsfile](https://github.com/cloudfoundry/cf-deployment/blob/master/operations/experimental/enable-backup-restore.yml)
at deploy time to ensure the backup and restore scripts are enabled. 
This will also deploy a backup restore VM, on which the 
[Backup and Restore SDK].(https://github.com/cloudfoundry-incubator/backup-and-restore-sdk-release) is deployed.


1. Spin up a Cloud Foundry deployment.
    * CF on BOSH Lite is supported.
    * [cf-deployment](https://github.com/cloudfoundry/cf-deployment) is supported. Ensure you apply the [backup-restore opsfile](https://github.com/cloudfoundry/cf-deployment/blob/master/operations/experimental/enable-backup-restore.yml) at deploy time to ensure the backup and restore scripts are enabled. This will also deploy a backup restore VM, on which the [Backup and Restore SDK].(https://github.com/cloudfoundry-incubator/backup-and-restore-sdk-release) is deployed.
1. Set environment variable `DEFAULT_TIMEOUT_MINS` if needed - timeout for commands run in the test. Defaults to 15 minutes.

1. Run `scripts/run_acceptance_tests_local.sh` with an integrations config set like this. Description of possible variables below.
    ```bash
    {
      "cf_api_url": "api.bosh-lite.com",
      "cf_deployment_name": "cf",
      "cf_admin_username": "admin",
      "cf_admin_password": "admin",
      "bosh_environment": "",
      "bosh_client": "",
      "bosh_client_secret": "",
      "bosh_ca_cert": "",
      "bosh_gw_host": "",
      "bosh_gw_user": "",
      "bosh_gw_private_key_contents": ""
    }

    ```
1. If you wish to run DRATS against a director deployed with `bbl`, run `scripts/run_acceptance_tests_with_bbl_env.sh <path-to-bbl-state-dir>`.
    * Set `CF_VARS_STORE_PATH` to the path to the CF vars-store file.
    * Set `BOSH_CLI_NAME` to the name of the BOSH CLI executable on your machine if it isn't `bosh`.
    * The script will search for NFS broker credentials in the CF vars-store file and will skip the NFS test case if those credentials are not present.

### Integration Config Variables
  * `cf_deployment_name` - name of the Cloud Foundry deployment to backup and restore
  * `cf_api_url` - Cloud Foundry api url
  * `cf_admin_username` - Cloud Foundry api admin user
  * `cf_admin_password` - Cloud Foundry api admin password
  * `bosh_environment` - URL of BOSH Director which has deployed the above Cloud Foundries
  * `bosh_client` - BOSH Director username
  * `bosh_client_secret` - BOSH Director password
  * `bosh_ca_cert` - BOSH Director's CA cert content
  * `bosh_gw_host` - Gateway host to use for BOSH SSH connection
  * `bosh_gw_user` - Gateway user to use for BOSH SSH connection
  * `bosh_gw_private_key_contents` - Private key to use for BOSH SSH connection
  
  ####Optional Variables
  * `ssh_destination_cidr` - Default to "10.0.0.0/8"; change if your cf-deployment is deployed in a different internal network range
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
1. Backs up the `cf_deployment_name` Cloud Foundry deployment.
1. Calls `AfterBackup(common.Config)` on all provided TestCases.
1. Restores to the `cf_deployment_name` Cloud Foundry deployment.
1. Calls `AfterRestore(common.Config)` on all provided TestCases (to e.g. check the apps pushed are present in the restored environment).
1. Calls `Cleanup(common.Config)` on all provided TestCases (to e.g. clean up the apps from the backup environment).

## Extending DRATS

DRATS runs a collection of test cases against a Cloud Foundry deployment.

Test cases should be used for checking that cf components' data has been backed up and restored correctly – e.g. if your release backs up a table in a database, that the table can be altered and is then restored to its original state.

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
