package testcases

import (
	"path"
	"time"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CfAppTestCase struct {
	uniqueTestID string
	appName      string
	envVarValue  string
	name         string
}

func NewCfAppTestCase() *CfAppTestCase {
	id := RandomStringNumber()
	return &CfAppTestCase{
		uniqueTestID: id,
		appName:      "test_app_" + id,
		envVarValue:  "winnebago" + id,
		name:         "cf-app",
	}
}

func (tc *CfAppTestCase) Name() string {
	return tc.name
}

func (tc *CfAppTestCase) BeforeBackup(config Config) {
	By("creating new orgs and spaces")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.ApiUrl)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	var testAppFixturePath = path.Join(CurrentTestDir(), "/../fixtures/test_app/")
	RunCommandSuccessfully("cf push " + tc.appName + " -p " + testAppFixturePath)
	RunCommandSuccessfully("cf set-env " + tc.appName + " MY_SPECIAL_VAR " + tc.envVarValue)
}

func (tc *CfAppTestCase) AfterBackup(config Config) {
	tc.deletePushedApps(config)
}

func (tc *CfAppTestCase) AfterRestore(config Config) {
	By("finding credentials for the deployment to restore")

	By("verify orgs and spaces have been re-created")
	RunCommandSuccessfully("cf org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf space acceptance-test-space-" + tc.uniqueTestID)

	By("verifying apps are back")
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	url := GetAppUrl(tc.appName)

	Eventually(StatusCode("https://"+url), 5*time.Minute, 5*time.Second).Should(Equal(200))
	Expect(string(RunCommandSuccessfully("cf env " + tc.appName).Out.Contents())).To(MatchRegexp("winnebago" + tc.uniqueTestID))
}

func (tc *CfAppTestCase) Cleanup(config Config) {
	tc.deletePushedApps(config)
}

func (tc *CfAppTestCase) deletePushedApps(config Config) {
	By("cleaning up orgs and spaces")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
