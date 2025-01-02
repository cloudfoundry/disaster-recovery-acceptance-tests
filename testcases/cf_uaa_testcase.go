package testcases

import (
	. "github.com/cloudfoundry/disaster-recovery-acceptance-tests/runner"

	. "github.com/onsi/ginkgo/v2"
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
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully("cf auth", username, password)
}

func (tc *CfUaaTestCase) Name() string {
	return tc.name
}

func (tc *CfUaaTestCase) CheckDeployment(config Config) {
}

func (tc *CfUaaTestCase) BeforeBackup(config Config) {
	By("creating a test user and logging in")
	login(config, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-user ", tc.testUser, tc.testPassword)
	login(config, tc.testUser, tc.testPassword)
	RunCommandSuccessfully("cf logout")
}

func (tc *CfUaaTestCase) AfterBackup(config Config) {
	By("deleting the test user")
	login(config, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf delete-user ", tc.testUser, "-f")
	RunCommandSuccessfully("cf logout")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)

	By("verifying that the test user cannot authenticate")
	result := RunCommand("cf auth", tc.testUser, tc.testPassword)
	Expect(result.ExitCode()).To(Equal(1))
}

func (tc *CfUaaTestCase) EnsureAfterSelectiveRestore(config Config) {}

func (tc *CfUaaTestCase) AfterRestore(config Config) {
	By("verifying that the test user can login again")
	login(config, tc.testUser, tc.testPassword)
}

func (tc *CfUaaTestCase) Cleanup(config Config) {
	By("deleting the test user")
	login(config, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf delete-user ", tc.testUser, "-f")
}
