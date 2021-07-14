package testcases

import (
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CfUaaTestCase struct {
	testUser     string
	testPassword string
	name         string
}

func NewCfUaaTestCase() *CfUaaTestCase {
	randomString := RandomStringNumber()
	testUser := "uaa-test-user-" + randomString
	testPassword := "uaa-test-password-" + randomString
	return &CfUaaTestCase{testUser: testUser, testPassword: testPassword, name: "cf-uaa"}
}

func login(config Config, username, password string) {
	RunCommandSuccessfully(CF_CLI+" api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully(CF_CLI+" auth", username, password)
}

func (tc *CfUaaTestCase) Name() string {
	return tc.name
}

func (tc *CfUaaTestCase) CheckDeployment(config Config) {
}

func (tc *CfUaaTestCase) BeforeBackup(config Config) {
	By("we create a user and can login")
	login(config, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully(CF_CLI+" create-user ", tc.testUser, tc.testPassword)
	login(config, tc.testUser, tc.testPassword)
	RunCommandSuccessfully(CF_CLI + " logout")
}

func (tc *CfUaaTestCase) AfterBackup(config Config) {
	By("we delete the user and verify")
	login(config, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully(CF_CLI+" delete-user ", tc.testUser, "-f")
	RunCommandSuccessfully(CF_CLI + " logout")
	RunCommandSuccessfully(CF_CLI+" api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)

	By("user has been deleted. authentication should fail")
	result := RunCommand(CF_CLI+" auth", tc.testUser, tc.testPassword)
	Expect(result.ExitCode()).To(Equal(1))
}

func (tc *CfUaaTestCase) EnsureAfterSelectiveRestore(config Config) {}

func (tc *CfUaaTestCase) AfterRestore(config Config) {
	By("we can login again")
	login(config, tc.testUser, tc.testPassword)
}

func (tc *CfUaaTestCase) Cleanup(config Config) {
	By("We delete the user")
	login(config, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully(CF_CLI+" delete-user ", tc.testUser, "-f")
}
