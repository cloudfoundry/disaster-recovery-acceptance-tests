# disaster-recovery-acceptance-tests (DRATS)

Tests if Cloud Foundry can be backed up and restored. The tests will back up `DEPLOYMENT_TO_BACKUP` and restore to `DEPLOYMENT_TO_RESTORE`. The two environments fed in should be identical to each other, i.e. they should have the same set of credentials in their deployment manifests. You can also run the tests on a single environment, by setting `DEPLOYMENT_TO_BACKUP` and `DEPLOYMENT_TO_RESTORE` to the same environment.

## Running DRATS

1. Spin up a Cloud Foundry deployment.
1. Run `scripts/run_acceptance_tests.sh` with the following environment variables set:
  * `DEPLOYMENT_TO_BACKUP` - name of the Cloud Foundry deployment to backup
  * `DEPLOYMENT_TO_RESTORE` - name of the Cloud Foundry deployment to restore
  * `BOSH_URL` - URL of BOSH Director which has deployed the above Cloud Foundries
  * `BOSH_CLIENT` - BOSH Director username
  * `BOSH_CLIENT_SECRET` - BOSH Director password
  * `BOSH_CERT_PATH` - path to BOSH Director's CA cert
  * `BBR_BUILD_PATH` - path to BBR binary

Currently it is recommended to have DRATS back up from and restore to the same environment.

## Test Structure

The system tests do the following:

1. Calls `BeforeBackup(common.Config)` on all provided TestCases (to e.g. push unique apps to the environment to be backed up).
1. Backs up the `DEPLOYMENT_TO_BACKUP` Cloud Foundry deployment.
1. Calls `AfterBackup(common.Config)` on all provided TestCases.
1. Restores to the `DEPLOYMENT_TO_RESTORE` Cloud Foundry deployment.
1. Calls `AfterRestore(common.Config)` on all provided TestCases (to e.g. check the apps pushed are present in the restored environment).
1. Calls `Cleanup(common.Config)` on all provided TestCases (to e.g. clean up the apps from the backup environment).

## Extending DRATS

DRATS runs a collection of test cases against two Cloud Foundry deployments.

To add extra test cases, create a new TestCase that follows the [TestCase interface](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/blob/master/runner/testcase.go).

The methods that need to be implemented are `BeforeBackup(common.Config)`, `AfterBackup(common.Config)`, `AfterRestore(common.Config)` and `Cleanup(common.Config)`.

* `BeforeBackup(common.Config)` runs before the backup is taken, and should create state in the Cloud Foundry deployment to be backed up (whose name is set to environment variable `DEPLOYMENT_TO_BACKUP`).
* `AfterBackup(common.Config)` runs after the backup is complete but before the restore is started. If were monitoring e.g. app uptime during the backup you could use this step to stop monitoring knowing that backup definitely finished.
* `AfterRestore(common.Config)` runs after the restore is complete, and should assert that the state in the restored Cloud Foundry deployment (whose name is set to environment variable `DEPLOYMENT_TO_RESTORE`) matches that created in `BeforeBackup(common.Config)`.
* `Cleanup(common.Config)` should clean up the state created in the Cloud Foundry deployment to be backed up.

`common.Config` contains the config for the BOSH Director and for the CF deployments to backup and restore.

1. Create a new TestCase in `test_cases`
1. In `testcases/testcase_helper.go`, initialise the TestCase and add it to the slice returned by `OpenSourceTestCases()`
