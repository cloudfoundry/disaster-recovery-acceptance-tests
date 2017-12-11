package runner

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
)

type TestContext struct {
	WorkspaceDir    string
	BinaryPath      string
	CertificatePath string
}

func NewTestContext(uniqueTestID string, boshConfig BoshConfig) (*TestContext, error) {
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

	var err error
	testContext.CertificatePath, err = writeBoshCaCertToFile(testContext.WorkspaceDir, boshConfig.BoshCaCert)

	return &testContext, err
}

func (testContext *TestContext) Cleanup() {
	By("remove workspace directory")
	RunCommandSuccessfullyWithFailureMessage("removing workspace directory", "rm -rf", testContext.WorkspaceDir)
}

func writeBoshCaCertToFile(tmpDir, boshCaCert string) (string, error) {
	dir, err := ioutil.TempDir(tmpDir, "drats")
	if err != nil {
		return "", err
	}

	boshCaCertFile, err := ioutil.TempFile(dir, "boshca")
	if err != nil {
		return "", err
	}

	_, err = boshCaCertFile.WriteString(boshCaCert)

	return boshCaCertFile.Name(), err
}
