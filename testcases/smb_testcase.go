package testcases

import (
	"fmt"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SMBTestCase struct {
	uniqueTestID string
	instanceName string
}

func NewSMBTestCases() *SMBTestCase {
	id := RandomStringNumber()
	name := fmt.Sprintf("service-instance-%s", id)
	return &SMBTestCase{uniqueTestID: id, instanceName: name}
}

func (tc *SMBTestCase) Name() string {
	return "cf-smbbroker"
}

func (tc *SMBTestCase) CheckDeployment(config Config) {
	By("checking if the SMB service is registered")
	RunCommandAndRetry(CF_CLI+" api --skip-ssl-validation", 3, config.CloudFoundryConfig.APIURL)
	RunCommandAndRetry(CF_CLI+" auth", 3, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfullyWithFailureMessage(
		tc.Name()+" test case cannot be run: SMB service is not registered",
		CF_CLI+" service-access -e "+config.CloudFoundryConfig.SMBServiceName,
	)
}

func (tc *SMBTestCase) BeforeBackup(config Config) {
	By("checking the service name and plane name are provided")
	Expect(config.CloudFoundryConfig.SMBServiceName).NotTo(BeEmpty(), "required config SMB service name not set")
	Expect(config.CloudFoundryConfig.SMBPlanName).NotTo(BeEmpty(), "required config SMB plan name not set")

	By("creating an SMB service broker and service instance")
	RunCommandSuccessfully(CF_CLI+" api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully(CF_CLI+" login --skip-ssl-validation -a", config.CloudFoundryConfig.APIURL,
		"-u", config.CloudFoundryConfig.AdminUsername, "-p", config.CloudFoundryConfig.AdminPassword)
	orgName := "acceptance-test-org-" + tc.uniqueTestID
	spaceName := "acceptance-test-space-" + tc.uniqueTestID
	RunCommandSuccessfully(CF_CLI + " create-org " + orgName)
	RunCommandSuccessfully(CF_CLI + " create-space " + spaceName + " -o " + orgName)
	RunCommandSuccessfully(CF_CLI + " target -o " + orgName + " -s " + spaceName)
	RunCommandSuccessfully(CF_CLI + " push dratsApp --docker-image docker/httpd --no-start --random-route")

	if config.CloudFoundryConfig.SMBCreateServiceBroker {
		RunCommandSuccessfully(CF_CLI + " create-service-broker " + "smbbroker-drats-" + tc.uniqueTestID + " " +
			config.CloudFoundryConfig.SMBBrokerUser + " " + config.CloudFoundryConfig.SMBBrokerPassword + " " +
			config.CloudFoundryConfig.SMBBrokerURL)
	}

	RunCommandSuccessfully(CF_CLI + " enable-service-access " + config.CloudFoundryConfig.SMBServiceName + " -o " + orgName)
	RunCommandSuccessfully(CF_CLI + " create-service " + config.CloudFoundryConfig.SMBServiceName + " " +
		config.CloudFoundryConfig.SMBPlanName + " " + tc.instanceName + " -c " +
		`'{"share":"//someserver.someplace.com/someshare"}'`)
}

func (tc *SMBTestCase) AfterBackup(config Config) {
	By("deleting the SMB service instance after backup")
	RunCommandSuccessfully(CF_CLI + " delete-service " + tc.instanceName + " -f")
}
func (tc *SMBTestCase) EnsureAfterSelectiveRestore(config Config) {
	By("repushing apps if restoring from a selective restore")
	RunCommandSuccessfully(CF_CLI + " push dratsApp --docker-image docker/httpd --no-start --random-route")
}

func (tc *SMBTestCase) AfterRestore(config Config) {
	By("re-binding the SMB service instance after restore")
	RunCommandSuccessfully(CF_CLI + " bind-service dratsApp " + tc.instanceName + ` -c '{"username":"someuser","password":"somepass"}'`)
}

func (tc *SMBTestCase) Cleanup(config Config) {
	By("smb cleanup")
	RunCommandSuccessfully(CF_CLI + " delete-org -f acceptance-test-org-" + tc.uniqueTestID)

	if config.CloudFoundryConfig.SMBCreateServiceBroker {
		RunCommandSuccessfully(CF_CLI + " delete-service-broker -f " + "smbbroker-drats-" + tc.uniqueTestID)
	}
}
