package testcases

import (
	"fmt"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
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

func (tc *SMBTestCase) BeforeBackup(config Config) {
	By("creating an SMB service broker and service instance")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.ApiUrl)
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.CloudFoundryConfig.ApiUrl,
		"-u", config.CloudFoundryConfig.AdminUsername, "-p", config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID +
		" -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID +
		" -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf push dratsApp --docker-image docker/httpd --no-start --random-route")

	if config.CloudFoundryConfig.SMBBrokerUser != "" {
		RunCommandSuccessfully("cf create-service-broker " + "smbbroker-drats-" + tc.uniqueTestID + " " +
			config.CloudFoundryConfig.SMBBrokerUser + " " + config.CloudFoundryConfig.SMBBrokerPassword + " " +
			config.CloudFoundryConfig.SMBBrokerUrl)
	}
	RunCommandSuccessfully("cf enable-service-access " + config.CloudFoundryConfig.SMBServiceName)
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
	if config.CloudFoundryConfig.SMBBrokerUser != "" {
		RunCommandSuccessfully("cf delete-service-broker -f " + "smbbroker-drats-" + tc.uniqueTestID)
	}
}

func (tc *SMBTestCase) deletePushedApps(config Config) {
}
