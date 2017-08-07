package common

import (
	"fmt"
	"io"
	"os"
	"strings"

	"bytes"

	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunCommandSuccessfully(cmd string, args ...string) *gexec.Session {
	session := runCommandWithStream(GinkgoWriter, GinkgoWriter, cmd, args...)
	Expect(session).To(gexec.Exit(0))
	return session
}

func RunCommand(cmd string, args ...string) *gexec.Session {
	return runCommandWithStream(GinkgoWriter, GinkgoWriter, cmd, args...)
}

func runCommandWithStream(stdout, stderr io.Writer, cmd string, args ...string) *gexec.Session {
	cmdToRunArgs := strings.Join(args, " ")
	cmdToRun := cmd + " " + cmdToRunArgs

	command := exec.Command("bash", "-c", cmdToRun)
	session, err := gexec.Start(command, stdout, stderr)

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit())
	return session
}

func DownloadManifest(deploymentName string, boshConfig BoshConfig) string {
	session := runCommandWithStream(bytes.NewBufferString(""), GinkgoWriter, join(
		BoshCommand(boshConfig),
		forDeployment(deploymentName),
		"manifest",
	))
	Eventually(session).Should(gexec.Exit(0))

	return string(session.Out.Contents())
}

func DeploymentToBackup() string {
	return MustHaveEnv("DEPLOYMENT_TO_BACKUP")
}

func DeploymentToRestore() string {
	return MustHaveEnv("DEPLOYMENT_TO_RESTORE")
}

func BoshCommand(boshConfig BoshConfig) string {
	return fmt.Sprintf("bosh-cli --non-interactive --environment=%s --ca-cert=%s --client=%s --client-secret=%s",
		boshConfig.BoshURL,
		boshConfig.BoshCertPath,
		boshConfig.BoshClient,
		boshConfig.BoshClientSecret,
	)
}

func forDeployment(deploymentName string) string {
	return fmt.Sprintf(
		"--deployment=%s",
		deploymentName,
	)
}

func MustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	if val == "" {
		panic(fmt.Sprintln("Env var %s not set", keyname))
	}
	return val
}

func join(args ...string) string {
	return strings.Join(args, " ")
}
