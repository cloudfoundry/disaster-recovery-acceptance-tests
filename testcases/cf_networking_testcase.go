package testcases

import (
	"fmt"
	"path"

	. "github.com/cloudfoundry/disaster-recovery-acceptance-tests/runner"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type CfNetworkingTestCase struct {
	uniqueTestID       string
	name               string
	testAppFixturePath string
	testAppName        string
}

func NewCfNetworkingTestCase() *CfNetworkingTestCase {
	id := RandomStringNumber()
	return &CfNetworkingTestCase{
		uniqueTestID:       id,
		name:               "cf-networking",
		testAppFixturePath: path.Join(CurrentTestDir(), "/../fixtures/test_app/"),
		testAppName:        fmt.Sprintf("test_app_%s", id),
	}
}

func (tc *CfNetworkingTestCase) Name() string {
	return tc.name
}

func (tc *CfNetworkingTestCase) CheckDeployment(config Config) {
}

func (tc *CfNetworkingTestCase) BeforeBackup(config Config) {
	By("creating a test org and space")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)

	By("setting up a test app")
	RunCommandSuccessfully("cf push " + tc.testAppName + " -p " + tc.testAppFixturePath)
	RunCommandSuccessfully(fmt.Sprintf("cf add-network-policy %s %s --port 8080 --protocol tcp", tc.testAppName, tc.testAppName))
}

func (tc *CfNetworkingTestCase) AfterBackup(config Config) {
	By("removing the network policy from the test app")
	testAppName := fmt.Sprintf("test_app_%s", tc.uniqueTestID)
	RunCommandSuccessfully(fmt.Sprintf("cf remove-network-policy %s %s --port 8080 --protocol tcp", testAppName, testAppName))
}

func (tc *CfNetworkingTestCase) EnsureAfterSelectiveRestore(config Config) {
	By("repushing the test app if restoring from a selective restore")
	RunCommandSuccessfully("cf push " + tc.testAppName + " -p " + tc.testAppFixturePath)
}

func (tc *CfNetworkingTestCase) AfterRestore(config Config) {
	By("checking that the network policy for the test app exists")
	session := RunCommand(fmt.Sprintf("cf network-policies"))
	testAppName := fmt.Sprintf("test_app_%s", tc.uniqueTestID)
	Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf(`%s.*%s.*tcp.*8080`, testAppName, testAppName)))
}

func (tc *CfNetworkingTestCase) Cleanup(config Config) {
	tc.deletePushedApps(config)
}

func (tc *CfNetworkingTestCase) deletePushedApps(config Config) {
	By("deleting the test org")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
