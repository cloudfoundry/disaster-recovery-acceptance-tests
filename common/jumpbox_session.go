package common

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type Session struct {
	WorkspaceDir    string
	BinaryPath      string
	CertificatePath string
}

func NewSession(uniqueTestID string) *Session {
	session := Session{}
	session.WorkspaceDir = "/tmp/backup_workspace" + uniqueTestID
	var bbrBuildPath = MustHaveEnv("BBR_BUILD_PATH")

	By("setting up the session")
	Eventually(RunCommandSuccessfully("sudo mkdir -p", session.WorkspaceDir, "&& sudo chmod 0777", session.WorkspaceDir)).Should(gexec.Exit(0))
	session.BinaryPath = bbrBuildPath
	session.CertificatePath = MustHaveEnv("BOSH_CERT_PATH")

	return &session
}

func (session *Session) Cleanup() {
	By("remove workspace directory")
	RunCommandSuccessfully("sudo rm -rf", session.WorkspaceDir)
}
