package testcases

import (
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CfUaaTestCase struct {
	uniqueTestID string
	testPassword string
}

func NewCfUaaTestCase() *CfUaaTestCase {
	id := RandomStringNumber()
	password := RandomStringNumber()
	return &CfUaaTestCase{uniqueTestID: id, testPassword: password}
}

func login(config Config, username, password string) {
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.DeploymentToBackup.ApiUrl, "-u", username, "-p", password)
}

func (tc *CfUaaTestCase) BeforeBackup(config Config) {
	By("we create a user and can login")
	login(config, config.DeploymentToBackup.AdminUsername, config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf create-user ", tc.uniqueTestID, tc.testPassword)
	login(config, tc.uniqueTestID, tc.testPassword)
	RunCommandSuccessfully("cf logout")
}

func (tc *CfUaaTestCase) AfterBackup(config Config) {
	By("we delete the user and verify")
	login(config, config.DeploymentToBackup.AdminUsername, config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf delete-user ", tc.uniqueTestID, "-f")
	RunCommandSuccessfully("cf logout")
	result := RunCommand("cf login --skip-ssl-validation -a", config.DeploymentToBackup.ApiUrl, "-u", tc.uniqueTestID, "-p", "password")
	Expect(result.ExitCode()).To(Equal(1))
}

func (tc *CfUaaTestCase) AfterRestore(config Config) {
	By("we can login again")
	login(config, tc.uniqueTestID, tc.testPassword)
}

func (tc *CfUaaTestCase) Cleanup(config Config) {
	By("We delete the user")
	login(config, config.DeploymentToBackup.AdminUsername, config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf delete-user ", tc.uniqueTestID, "-f")
}

