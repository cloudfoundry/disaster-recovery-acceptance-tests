package testcases

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/fixtures"
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"

	. "github.com/onsi/ginkgo"
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
		uniqueTestID: id,
		name:         "cf-networking",
		testAppName:  fmt.Sprintf("test_app_%s", id),
	}
}

func (tc *CfNetworkingTestCase) Name() string {
	return tc.name
}

func (tc *CfNetworkingTestCase) CheckDeployment(config Config) {
}

func (tc *CfNetworkingTestCase) BeforeBackup(config Config) {
	tmpDir, err := ioutil.TempDir("", "cf-app-test-case-fixtures")
	Expect(err).NotTo(HaveOccurred())
	err = fixtures.WriteFixturesToTemporaryDirectory(tmpDir, "test_app")
	Expect(err).NotTo(HaveOccurred())
	tc.testAppFixturePath = filepath.Join(tmpDir, "test_app")

	By("creating new orgs and spaces")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf push " + tc.testAppName + " -p " + tc.testAppFixturePath)
	RunCommandSuccessfully(fmt.Sprintf("cf add-network-policy %s --destination-app %s --port 8080 --protocol tcp", tc.testAppName, tc.testAppName))
}

func (tc *CfNetworkingTestCase) AfterBackup(config Config) {
	testAppName := fmt.Sprintf("test_app_%s", tc.uniqueTestID)
	RunCommandSuccessfully(fmt.Sprintf("cf remove-network-policy %s --destination-app %s --port 8080 --protocol tcp", testAppName, testAppName))
}

func (tc *CfNetworkingTestCase) EnsureAfterSelectiveRestore(config Config) {
	By("repushing apps if restoring from a selective restore")
	RunCommandSuccessfully("cf push " + tc.testAppName + " -p " + tc.testAppFixturePath)
}

func (tc *CfNetworkingTestCase) AfterRestore(config Config) {
	By("finding credentials for the deployment to restore")
	session := RunCommand(fmt.Sprintf("cf network-policies"))
	testAppName := fmt.Sprintf("test_app_%s", tc.uniqueTestID)
	Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf(`%s.*%s.*tcp.*8080`, testAppName, testAppName)))
}

func (tc *CfNetworkingTestCase) Cleanup(config Config) {
	tc.deletePushedApps(config)

	err := os.RemoveAll(tc.testAppFixturePath)
	Expect(err).NotTo(HaveOccurred())
}

func (tc *CfNetworkingTestCase) deletePushedApps(config Config) {
	By("cleaning up orgs and spaces")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
