package testcases

import (
	"fmt"
	"path"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type CfNetworkingTestCase struct {
	uniqueTestID string
}

func NewCfNetworkingTestCase() *CfNetworkingTestCase {
	id := RandomStringNumber()
	return &CfNetworkingTestCase{uniqueTestID: id}
}

func (tc *CfNetworkingTestCase) BeforeBackup(config Config) {
	By("creating new orgs and spaces")
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.DeploymentToBackup.ApiUrl, "-u", config.DeploymentToBackup.AdminUsername, "-p", config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	var testAppFixturePath = path.Join(CurrentTestDir(), "/../fixtures/test_app/")
	testAppName := fmt.Sprintf("test_app_%s", tc.uniqueTestID)
	RunCommandSuccessfully("cf push " + testAppName + " -p " + testAppFixturePath)
	RunCommandSuccessfully(fmt.Sprintf("cf add-network-policy %s --destination-app %s --port 8080 --protocol tcp", testAppName, testAppName))
}

func (tc *CfNetworkingTestCase) AfterBackup(config Config) {
	testAppName := fmt.Sprintf("test_app_%s", tc.uniqueTestID)
	RunCommandSuccessfully(fmt.Sprintf("cf remove-network-policy %s --destination-app %s --port 8080 --protocol tcp", testAppName, testAppName))
}

func (tc *CfNetworkingTestCase) AfterRestore(config Config) {
	By("finding credentials for the deployment to restore")
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.DeploymentToRestore.ApiUrl, "-u", config.DeploymentToRestore.AdminUsername, "-p", config.DeploymentToRestore.AdminPassword)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	session := RunCommand(fmt.Sprintf("cf network-policies"))
	testAppName := fmt.Sprintf("test_app_%s", tc.uniqueTestID)
	Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf(`%s.*%s.*tcp.*8080`, testAppName, testAppName)))
}

func (tc *CfNetworkingTestCase) Cleanup(config Config) {
	tc.deletePushedApps(config)
}

func (tc *CfNetworkingTestCase) deletePushedApps(config Config) {
	By("cleaning up orgs and spaces")
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.DeploymentToBackup.ApiUrl, "-u", config.DeploymentToBackup.AdminUsername, "-p", config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf delete-space -f acceptance-test-space-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
