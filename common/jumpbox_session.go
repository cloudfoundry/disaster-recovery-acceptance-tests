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
	bbrFilename := getTar(bbrBuildPath)
	CopyOnJumpbox(filepath.Join(bbrBuildPath, bbrFilename), access_vm+session.WorkspaceDir)
	Eventually(RunOnJumpboxAsVcap("tar -C", session.WorkspaceDir, "-xvf", filepath.Join(session.WorkspaceDir, bbrFilename))).Should(gexec.Exit(0))
	session.BinaryPath = filepath.Join(session.WorkspaceDir, "releases", "bbr")

	By("copying the Director certificate to the jumpbox")
	CopyOnJumpbox(MustHaveEnv("BOSH_CERT_PATH"), access_vm+session.WorkspaceDir+"/bosh.crt")
	session.CertificatePath = filepath.Join(session.WorkspaceDir, "bosh.crt")

	return &session
}

func (jumpBoxSession *JumpBoxSession) Cleanup() {
	By("remove workspace directory")
	RunOnJumpbox("sudo rm -rf", jumpBoxSession.WorkspaceDir)
}

func getTar(path string) string {
	glob := filepath.Join(path, "*.tar")
	matches, err := filepath.Glob(glob)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(matches)).To(Equal(1), "There should be only one tar file present in the BBR binary path")
	return filepath.Base(matches[0])
}
