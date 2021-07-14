package testcases

import (
	"path"
	"time"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
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
	By("creating new orgs and spaces")
	RunCommandSuccessfully(CF_CLI+" api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully(CF_CLI+" auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully(CF_CLI + " create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " push " + tc.deletedAppName + " -p " + tc.testAppFixturePath)
	RunCommandSuccessfully(CF_CLI + " push " + tc.stoppedAppName + " -p " + tc.testAppFixturePath)
	RunCommandSuccessfully(CF_CLI + " push " + tc.runningAppName + " -p " + tc.testAppFixturePath)

	RunCommandSuccessfully(CF_CLI + " set-env " + tc.deletedAppName + " MY_SPECIAL_VAR " + tc.envVarValue)
	RunCommandSuccessfully(CF_CLI + " set-env " + tc.stoppedAppName + " MY_STOPPED_SPECIAL_VAR " + tc.envVarValue)
	RunCommandSuccessfully(CF_CLI + " set-env " + tc.runningAppName + " MY_RUNNING_SPECIAL_VAR " + tc.envVarValue)

	RunCommandSuccessfully(CF_CLI + " stop " + tc.stoppedAppName)
}

func (tc *CfAppTestCase) AfterBackup(config Config) {
	RunCommandSuccessfully(CF_CLI + " delete " + tc.deletedAppName)
	RunCommandSuccessfully(CF_CLI + " start " + tc.stoppedAppName)
	RunCommandSuccessfully(CF_CLI + " stop " + tc.runningAppName)
}

func (tc *CfAppTestCase) EnsureAfterSelectiveRestore(config Config) {
	By("repushing apps if restoring from a selective restore")
	RunCommandSuccessfully(CF_CLI + " target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " push " + tc.deletedAppName + " -p " + tc.testAppFixturePath)

	tc.failedToRestoreDroplets = true
}

func (tc *CfAppTestCase) AfterRestore(config Config) {
	By("finding credentials for the deployment to restore")

	By("verify orgs and spaces have been re-created")
	RunCommandSuccessfully(CF_CLI + " org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " target -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " space acceptance-test-space-" + tc.uniqueTestID)

	By("verifying apps are back")
	RunCommandSuccessfully(CF_CLI + " target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)

	deletedAppUrl := GetAppURL(tc.deletedAppName)
	Eventually(StatusCode("https://"+deletedAppUrl), 5*time.Minute, 5*time.Second).Should(Equal(200))
	Expect(string(RunCommandSuccessfully(CF_CLI + " env " + tc.deletedAppName).Out.Contents())).To(MatchRegexp("winnebago" + tc.uniqueTestID))

	stoppedAppUrl := GetAppURL(tc.stoppedAppName)
	Eventually(StatusCode("https://"+stoppedAppUrl), 5*time.Minute, 5*time.Second).Should(Equal(404))
	Expect(GetRequestedState(tc.stoppedAppName)).To(Equal("stopped"))
	Expect(string(RunCommandSuccessfully(CF_CLI + " env " + tc.stoppedAppName).Out.Contents())).To(MatchRegexp("winnebago" + tc.uniqueTestID))

	// when a selective restore occurs we know this app won't be running as the droplet won't exist, so lets not assert it.
	if !tc.failedToRestoreDroplets {
		runningAppUrl := GetAppURL(tc.runningAppName)
		Eventually(StatusCode("https://"+runningAppUrl), 5*time.Minute, 5*time.Second).Should(Equal(200))
		Expect(string(RunCommandSuccessfully(CF_CLI + " env " + tc.runningAppName).Out.Contents())).To(MatchRegexp("winnebago" + tc.uniqueTestID))
	}
}

func (tc *CfAppTestCase) Cleanup(config Config) {
	tc.deletePushedApps(config)
}

func (tc *CfAppTestCase) deletePushedApps(config Config) {
	By("cleaning up orgs and spaces")
	RunCommandSuccessfully(CF_CLI + " delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
