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
	By("checking if the smbbroker app is present")
	RunCommandAndRetry("cf api --skip-ssl-validation", 3, config.CloudFoundryConfig.ApiUrl)
	RunCommandAndRetry("cf auth", 3, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfullyWithFailureMessage(tc.Name()+" test case cannot be run: space smb-broker-space is not present", "cf target -o system -s smb-broker-space")
	RunCommandSuccessfullyWithFailureMessage(tc.Name()+" test case cannot be run: app smbbroker is not present", "cf app smbbroker")
}

func (tc *SMBTestCase) BeforeBackup(config Config) {
	By("checking the service name and plane name are provided")
	Expect(config.CloudFoundryConfig.SMBServiceName).NotTo(BeEmpty(), "required config SMB service name not set")
	Expect(config.CloudFoundryConfig.SMBPlanName).NotTo(BeEmpty(), "required config SMB plan name not set")

	By("creating an SMB service broker and service instance")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.ApiUrl)
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.CloudFoundryConfig.ApiUrl,
		"-u", config.CloudFoundryConfig.AdminUsername, "-p", config.CloudFoundryConfig.AdminPassword)
	orgName := "acceptance-test-org-" + tc.uniqueTestID
	spaceName := "acceptance-test-space-" + tc.uniqueTestID
	RunCommandSuccessfully("cf create-org " + orgName)
	RunCommandSuccessfully("cf create-space " + spaceName + " -o " + orgName)
	RunCommandSuccessfully("cf target -o " + orgName + " -s " + spaceName)
	RunCommandSuccessfully("cf push dratsApp --docker-image docker/httpd --no-start --random-route")

	if config.CloudFoundryConfig.SMBCreateServiceBroker {
		RunCommandSuccessfully("cf create-service-broker " + "smbbroker-drats-" + tc.uniqueTestID + " " +
			config.CloudFoundryConfig.SMBBrokerUser + " " + config.CloudFoundryConfig.SMBBrokerPassword + " " +
			config.CloudFoundryConfig.SMBBrokerUrl)
	}

	RunCommandSuccessfully("cf enable-service-access " + config.CloudFoundryConfig.SMBServiceName + " -o " + orgName)
	RunCommandSuccessfully("cf create-service " + config.CloudFoundryConfig.SMBServiceName + " " +
		config.CloudFoundryConfig.SMBPlanName + " " + tc.instanceName + " -c " +
		`'{"share":"//someserver.someplace.com/someshare"}'`)
}

func (tc *SMBTestCase) AfterBackup(config Config) {
	By("deleting the SMB service instance after backup")
	RunCommandSuccessfully("cf delete-service " + tc.instanceName + " -f")
}

func (tc *SMBTestCase) AfterRestore(config Config) {
	By("re-binding the SMB service instance after restore")
	RunCommandSuccessfully("cf bind-service dratsApp " + tc.instanceName + ` -c '{"username":"someuser","password":"somepass"}'`)
}

func (tc *SMBTestCase) Cleanup(config Config) {
	By("smb cleanup")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)

	if config.CloudFoundryConfig.SMBCreateServiceBroker {
		RunCommandSuccessfully("cf delete-service-broker -f " + "smbbroker-drats-" + tc.uniqueTestID)
	}
}
