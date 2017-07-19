# disaster-recovery-acceptance-tests (DRATS)

## Running DRATS

1. Spin up a Cloud Foundry deployment.
1. Deploy a jumpbox deployment called `integration-jump-box` containing a single VM, `jumpbox`.
1. Run `ci/scripts/acceptance.sh` with the following environment variables set:
  * `DEPLOYMENT_TO_BACKUP` - name of the Cloud Foundry deployment
  * `DEPLOYMENT_TO_RESTORE` - name of the Cloud Foundry deployment
  * `BOSH_URL` - URL of BOSH Director which has deployed the above Cloud Foundries
  * `BOSH_CLIENT` - BOSH Director username
  * `BOSH_CLIENT_SECRET` - BOSH Director password
  * `BOSH_CERT_PATH` - path to BOSH Director's CA cert
  * `BOSH_GATEWAY_USER` - BOSH SSH client username
  * `BOSH_GATEWAY_HOST` - BOSH SSH client hostname
  * `BOSH_GATEWAY_KEY` - path to BOSH SSH client private key
  * `BBR_BUILD_PATH` - path to BBR binary

Currently, DRATS backs up from and restores to the same environment.

## Test Structure

The system tests do the following:

1. Starts a session on the jumpbox VM (creates a workspace directory, copies over the BOSH Director CA cert and key and the BBR binary)
1. Calls `PopulateState()` on all provided TestCases (to e.g. push unique apps to the environment to be backed up).
1. Backs up the `DEPLOYMENT_TO_BACKUP` Cloud Foundry deployment.
1. Restores to the `DEPLOYMENT_TO_RESTORE` Cloud Foundry deployment.
1. Calls `CheckState()` on all provided TestCases (to e.g. check the apps pushed in (2) are present in the restored environment).
1. Calls `Cleanup()` on all provided TestCases (to e.g. clean up the apps from the backup environment).

## Extending DRATS

DRATS runs a collection of test cases against two Cloud Foundry deployments.

To add extra test cases, create a new TestCase that follows the [TestCase interface](https://github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/blob/master/acceptance/testcases/test_case.go).

The methods that need to be implemented are `PopulateState()`, `CheckState()` and `Cleanup()`.

* `PopulateState()` should create some state in the Cloud Foundry deployment to be backed up (whose name is set to environment variable `DEPLOYMENT_TO_BACKUP`).
* `CheckState()` should assert that the state in the restored Cloud Foundry deployment (whose name is set to environment variable `DEPLOYMENT_TO_RESTORE`) matches that created by `PopulateState()`.
* `Cleanup()` should clean up the state created in the Cloud Foundry deployment to be backed up.

1. Create a new TestCase in acceptance/testcases
1. In `acceptance/acceptance_suite_test.go`, initialise the TestCase and add it to the `testCases` slice