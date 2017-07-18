package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunOnRemote(deploymentName, instanceName string, cmd ...string) *gexec.Session {
	return RunCommand(
		join(
			BoshCommand(),
			forDeployment(deploymentName),
			getSSHCommand(instanceName, "0"),
		),
		join(cmd...),
	)
}

func RunOnJumpboxAsVcap(cmd ...string) *gexec.Session {
	return RunOnJumpbox("sudo", "su", "vcap", "-c", fmt.Sprintf("'%s'", join(cmd...)))
}

func RunOnJumpbox(cmd ...string) *gexec.Session {
	return RunOnRemote(JumpboxDeployment(), "jumpbox", cmd...)
}

func CopyOnJumpbox(source, destination string) {
	RunCommandSuccessfully(
		join(
			BoshCommand(),
			forDeployment(JumpboxDeployment()),
			getSCPCommand(),
			source, destination,
		),
	)
}

func RunCommandSuccessfully(cmd string, args ...string) *gexec.Session {
	session := runCommandWithStream(GinkgoWriter, GinkgoWriter, cmd, args...)
	Expect(session).To(gexec.Exit(0))
	return session
}

func RunCommand(cmd string, args ...string) *gexec.Session {
	return runCommandWithStream(GinkgoWriter, GinkgoWriter, cmd, args...)
}

func runCommandWithStream(stdout, stderr io.Writer, cmd string, args ...string) *gexec.Session {
	cmdParts := strings.Split(cmd, " ")
	commandPath := cmdParts[0]
	combinedArgs := append(cmdParts[1:], args...)
	command := exec.Command(commandPath, combinedArgs...)

	session, err := gexec.Start(command, stdout, stderr)

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit())
	return session
}

func DownloadManifest(deploymentName string) string {
	session := runCommandWithStream(bytes.NewBufferString(""), GinkgoWriter, join(
		BoshCommand(),
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

func BoshCommand() string {
	return fmt.Sprintf("bosh-cli --non-interactive --environment=%s --ca-cert=%s --client=%s --client-secret=%s",
		MustHaveEnv("BOSH_URL"),
		MustHaveEnv("BOSH_CERT_PATH"),
		MustHaveEnv("BOSH_CLIENT"),
		MustHaveEnv("BOSH_CLIENT_SECRET"),
	)
}

func forDeployment(deploymentName string) string {
	return fmt.Sprintf(
		"--deployment=%s",
		deploymentName,
	)
}

func JumpboxDeployment() string {
	return "integration-jump-box"
}

func getSSHCommand(instanceName, instanceIndex string) string {
	return fmt.Sprintf(
		"ssh --gw-user=%s --gw-host=%s --gw-private-key=%s %s/%s",
		MustHaveEnv("BOSH_GATEWAY_USER"),
		MustHaveEnv("BOSH_GATEWAY_HOST"),
		MustHaveEnv("BOSH_GATEWAY_KEY"),
		instanceName,
		instanceIndex,
	)
}

func getSCPCommand() string {
	return fmt.Sprintf(
		"scp --gw-user=%s --gw-host=%s --gw-private-key=%s",
		MustHaveEnv("BOSH_GATEWAY_USER"),
		MustHaveEnv("BOSH_GATEWAY_HOST"),
		MustHaveEnv("BOSH_GATEWAY_KEY"),
	)
}

func MustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	Expect(val).NotTo(BeEmpty(), "Need "+keyname+" for the test")
	return val
}

func join(args ...string) string {
	return strings.Join(args, " ")
}

func GetDebugFlag() string {
	debugFlag := ""
	if os.Getenv("DEBUG") == "true" {
		debugFlag = "--debug"
	}

	return debugFlag
}
