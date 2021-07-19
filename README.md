# disaster-recovery-acceptance-tests (DRATs): cf6-compatible branch

Tests if Cloud Foundry (CF) can be backed up and restored. The tests will back up from and restore to `CF_DEPLOYMENT_NAME`.
Specifically, DRATs adds state to a Cloud Foundry deployment by testing backup and restore during a CF operation such as pushing an app. DRATs backs up the deployment, restores from the backup, and asserts that the state is present after restore.

## cf6-compatible branch
This branch of DRATS uses the V6 version of the CF CLI. It is
therefore compatible with older CF deployments which cannot talk to
the V7 version of the CF CLI. The `main` branch of DRATS at the time
of writing uses the V7 version of the CF CLI.

This branch can be deleted or deprecated when all stakeholders have
moved to CF versions that support the V7 version of the CLI


## Prerequisites
1. [Install `go`](https://golang.org/)
1. [Install `ginkgo`](https://github.com/onsi/ginkgo)

## Running DRATs with Environment Variables

1. Spin up a Cloud Foundry deployment**.
    * The [backup-restore ops file](https://github.com/cloudfoundry/cf-deployment/blob/master/operations/backup-and-restore/enable-backup-restore.yml) must be used during deployment. This will also deploy a backup restore VM.
    * Other [backup-restore ops files](https://github.com/cloudfoundry/cf-deployment/blob/master/operations/backup-and-restore) will be needed depending on what test cases you want enabled.
2. Run DRATs against it using one of the following:
    * [Using an integration config](docs/testing_with_config.md) - **Recommended**
    * [Using bbl](docs/testing_with_bbl.md)
    * [Using environment variables](docs/testing_with_env_vars.md) - **Deprecated**

**Cloud Foundry on BOSH Lite is supported.


## Test Structure

The system tests do the following:

1. Calls `CheckDeployment(common.Config)` on all provided TestCases (to e.g. check if errand-pushed apps are present).
1. Sets up a temporary local working directory for storing the backup artifact, and CF_HOME directories for all the test cases.
1. Calls `BeforeBackup(common.Config)` on all provided TestCases (to e.g. push unique apps to the environment to be backed up).
1. Backs up the `CF_DEPLOYMENT_NAME` Cloud Foundry deployment.
1. Calls `AfterBackup(common.Config)` on all provided TestCases.
1. Restores to the `CF_DEPLOYMENT_NAME` Cloud Foundry deployment.
1. Calls `AfterRestore(common.Config)` on all provided TestCases (to e.g. check the apps pushed are present in the restored environment).
1. Calls `Cleanup(common.Config)` on all provided TestCases (to e.g. clean up the apps from the backup environment). It will do this even if an error or failure occurred in a previous step
1. Cleans up the temporary directories created in the setup
1. If an error occurred during a `bbr backup` command, DRATs runs `bbr backup-cleanup` to remove temporary bbr artifacts from your deployment (which would otherwise cause subsequent DRATs runs to fail)

## Extending DRATs

DRATs runs a collection of test cases against a Cloud Foundry deployment.

Test cases should be used for checking that CF components' data has been backed up and restored correctly – e.g. if your release backs up a table in a database, that the table can be altered and is then restored to its original state.

**Backup and restore of apps is covered by the existing CAPI test case.** No new test cases are needed for this – unless you're writing CAPI backup and restore scripts, app backup and restore can be assumed to work.

To add extra test cases, create a new TestCase that follows the [TestCase interface](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/blob/master/runner/testcase.go).

The methods that need to be implemented are:
* `CheckDeployment(common.Config)` runs first to ensure that the test case could possibly succeed in the deployment DRATs is running against. E.g. An errand-based test case could never succeed in a deployment where the errand has never been run.
* `BeforeBackup(common.Config)` runs before the backup is taken, and should create state in the Cloud Foundry deployment to be backed up.
* `AfterBackup(common.Config)` runs after the backup is complete but before the restore is started. If were monitoring e.g. app uptime during the backup you could use this step to stop monitoring knowing that backup definitely finished.
* `AfterRestore(common.Config)` runs after the restore is complete, and should assert that the state in the restored Cloud Foundry deployment matches that created in `BeforeBackup(common.Config)`.
* `Cleanup(common.Config)` should clean up the state created in the Cloud Foundry deployment to be backed up.

`common.Config` contains the config for the BOSH Director and for the CF deployments to backup and restore.

1. Create a new TestCase in `test_cases`
1. In `testcases/testcase_helper.go`, initialise the TestCase and add it to the slice returned by `OpenSourceTestCases()`

## Running DRATs in your CI

We provide tasks to run DRATs with your CI:
* A [`drats` task](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/tree/master/ci/drats) that reads in environment variables
* A [`drats-with-integration-config` task](https://github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/tree/master/ci/drats-with-integration-config) that read from an integration config.

Both DRATs tasks establish an SSH tunnel using [`sshuttle`](http://sshuttle.readthedocs.io) so that they can run from outside the network. Note the tasks will need to be run from a privileged container.
You can also find [our pipeline definition here](https://github.com/cloudfoundry-incubator/backup-and-restore-ci/blob/master/pipelines/drats/pipeline.yml)

## Debugging your DRATs run

DRATs runs multiple interwoven test cases (for app uptime and each of the components under test) so it can be a little tricky to work out what's gone wrong when there's an error or failure. Here are some tips on investigating DRATs failures - please PR in additions to this doc if you think of more tips that might help other teams!

1. The `bbr backup-cleanup` command runs if the test run errored during the `bbr backup` step. If you see an error in the `backup-cleanup` step, it's likely that a similar problem happened in the `backup` step which caused the original failure - scroll up to see.
1. The easiest way to see where the failure / error happened is to look for the nearest `STEP` statement in the logs
