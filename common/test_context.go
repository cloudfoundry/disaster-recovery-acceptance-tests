package common

import (
	. "github.com/onsi/ginkgo"
)

type TestContext struct {
	WorkspaceDir    string
	BinaryPath      string
	CertificatePath string
}

func NewTestContext(uniqueTestID string, boshConfig BoshConfig) *TestContext {
	testContext := TestContext{}
	testContext.WorkspaceDir = "/tmp/backup_workspace" + uniqueTestID
	var bbrBuildPath = MustHaveEnv("BBR_BUILD_PATH")

	By("setting up the test context")
	RunCommandSuccessfullyWithFailureMessage(
		"creating workspace directory",
		"mkdir -p", testContext.WorkspaceDir,
		"&& chmod 0777",
		testContext.WorkspaceDir)
	testContext.BinaryPath = bbrBuildPath
	testContext.CertificatePath = boshConfig.BoshCertPath

	return &testContext
}

func (testContext *TestContext) Cleanup() {
	By("remove workspace directory")
	RunCommandSuccessfullyWithFailureMessage("removing workspace directory", "rm -rf", testContext.WorkspaceDir)
}
