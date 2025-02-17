package testcases

import (
	"path"
	"time"

	. "github.com/cloudfoundry/disaster-recovery-acceptance-tests/runner"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type CfAppTestCase struct {
	uniqueTestID            string
	deletedAppName          string
	stoppedAppName          string
	runningAppName          string
	envVarValue             string
	name                    string
	testAppFixturePath      string
	failedToRestoreDroplets bool
}

func NewCfAppTestCase() *CfAppTestCase {
	id := RandomStringNumber()
	return &CfAppTestCase{
		uniqueTestID:       id,
		deletedAppName:     "test_app_" + id,
		stoppedAppName:     "stopped_test_app_" + id,
		runningAppName:     "running_test_app_" + id,
		envVarValue:        "winnebago" + id,
		name:               "cf-app",
		testAppFixturePath: path.Join(CurrentTestDir(), "/../fixtures/test_app/"),
	}
}

func (tc *CfAppTestCase) Name() string {
	return tc.name
}

func (tc *CfAppTestCase) CheckDeployment(config Config) {
}

func (tc *CfAppTestCase) BeforeBackup(config Config) {
	By("creating a test org and space")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)

	By("setting up some test apps")
	RunCommandSuccessfully("cf push " + tc.deletedAppName + " -p " + tc.testAppFixturePath)
	RunCommandSuccessfully("cf push " + tc.stoppedAppName + " -p " + tc.testAppFixturePath)
	RunCommandSuccessfully("cf push " + tc.runningAppName + " -p " + tc.testAppFixturePath)

	RunCommandSuccessfully("cf set-env " + tc.deletedAppName + " MY_SPECIAL_VAR " + tc.envVarValue)
	RunCommandSuccessfully("cf set-env " + tc.stoppedAppName + " MY_STOPPED_SPECIAL_VAR " + tc.envVarValue)
	RunCommandSuccessfully("cf set-env " + tc.runningAppName + " MY_RUNNING_SPECIAL_VAR " + tc.envVarValue)

	RunCommandSuccessfully("cf stop " + tc.stoppedAppName)
}

func (tc *CfAppTestCase) AfterBackup(config Config) {
	By("changing the state of the test apps")
	RunCommandSuccessfully("cf delete -f " + tc.deletedAppName)
	RunCommandSuccessfully("cf start " + tc.stoppedAppName)
	RunCommandSuccessfully("cf stop " + tc.runningAppName)
}

func (tc *CfAppTestCase) EnsureAfterSelectiveRestore(config Config) {
	By("repushing the deleted app if restoring from a selective restore")
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf push " + tc.deletedAppName + " -p " + tc.testAppFixturePath)

	tc.failedToRestoreDroplets = true
}

func (tc *CfAppTestCase) AfterRestore(config Config) {
	By("finding credentials for the deployment to restore")

	By("verify orgs and spaces have been re-created")
	RunCommandSuccessfully("cf org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf space acceptance-test-space-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)

	By("verifying that all apps are back")
	deletedAppUrl := GetAppURL(tc.deletedAppName)
	Eventually(StatusCode("https://"+deletedAppUrl), 5*time.Minute, 5*time.Second).Should(Equal(200))
	Expect(string(RunCommandSuccessfully("cf env " + tc.deletedAppName).Out.Contents())).To(MatchRegexp("winnebago" + tc.uniqueTestID))

	stoppedAppUrl := GetAppURL(tc.stoppedAppName)
	Eventually(StatusCode("https://"+stoppedAppUrl), 5*time.Minute, 5*time.Second).Should(Equal(404))
	Expect(GetRequestedState(tc.stoppedAppName)).To(Equal("stopped"))
	Expect(string(RunCommandSuccessfully("cf env " + tc.stoppedAppName).Out.Contents())).To(MatchRegexp("winnebago" + tc.uniqueTestID))

	// when a selective restore occurs we know this app won't be running as the droplet won't exist, so lets not assert it.
	if !tc.failedToRestoreDroplets {
		runningAppUrl := GetAppURL(tc.runningAppName)
		Eventually(StatusCode("https://"+runningAppUrl), 5*time.Minute, 5*time.Second).Should(Equal(200))
		Expect(string(RunCommandSuccessfully("cf env " + tc.runningAppName).Out.Contents())).To(MatchRegexp("winnebago" + tc.uniqueTestID))
	}
}

func (tc *CfAppTestCase) Cleanup(config Config) {
	tc.deletePushedApps(config)
}

func (tc *CfAppTestCase) deletePushedApps(config Config) {
	By("deleting the test org")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
