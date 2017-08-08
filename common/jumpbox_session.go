package common

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

//TODO think of the correct name for this

type Session struct {
	WorkspaceDir    string
	BinaryPath      string
	CertificatePath string
}

func NewSession(uniqueTestID string, boshConfig BoshConfig) *Session {
	session := Session{}
	session.WorkspaceDir = "/tmp/backup_workspace" + uniqueTestID
	var bbrBuildPath = MustHaveEnv("BBR_BUILD_PATH")

	By("setting up the session")
	Eventually(RunCommandSuccessfully("mkdir -p", session.WorkspaceDir, "&& chmod 0777", session.WorkspaceDir)).Should(gexec.Exit(0))
	session.BinaryPath = bbrBuildPath
	session.CertificatePath = boshConfig.BoshCertPath

	return &session
}

func (session *Session) Cleanup() {
	By("remove workspace directory")
	RunCommandSuccessfully("rm -rf", session.WorkspaceDir)
}
