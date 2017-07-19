package common

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type JumpBoxSession struct {
	WorkspaceDir    string
	BinaryPath      string
	CertificatePath string
}

func NewJumpBoxSession(uniqueTestID string) *JumpBoxSession {
	access_vm := "jumpbox/0:"

	session := JumpBoxSession{}
	session.WorkspaceDir = "/var/vcap/store/backup_workspace" + uniqueTestID
	var bbrBuildPath = MustHaveEnv("BBR_BUILD_PATH")

	By("setting up the jump box")
	Eventually(RunOnJumpboxAsVcap("sudo mkdir", session.WorkspaceDir, "&& sudo chown vcap:vcap", session.WorkspaceDir, "&& sudo chmod 0777", session.WorkspaceDir)).Should(gexec.Exit(0))

	By("copying the BBR binary to the jumpbox")
	CopyOnJumpbox(bbrBuildPath, access_vm+session.WorkspaceDir)
	session.BinaryPath = filepath.Join(session.WorkspaceDir, filepath.Base(bbrBuildPath))

	By("copying the Director certificate to the jumpbox")
	CopyOnJumpbox(MustHaveEnv("BOSH_CERT_PATH"), access_vm+session.WorkspaceDir+"/bosh.crt")
	session.CertificatePath = filepath.Join(session.WorkspaceDir, "bosh.crt")

	return &session
}

func (jumpBoxSession *JumpBoxSession) Cleanup() {
	By("remove workspace directory")
	RunOnJumpbox("sudo rm -rf", jumpBoxSession.WorkspaceDir)
}
