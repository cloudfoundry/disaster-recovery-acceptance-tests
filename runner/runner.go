package runner

import (
	"fmt"

	"time"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunDisasterRecoveryAcceptanceTests(configGetter ConfigGetter, testCases []TestCase) {
	var envsAreSame bool
	var uniqueTestID string
	var session *Session
	var config Config

	BeforeEach(func() {
		config = configGetter.FindConfig()

		SetDefaultEventuallyTimeout(30 * time.Minute)
		uniqueTestID = RandomStringNumber()
		session = NewSession(uniqueTestID, config.BoshConfig)
	})

	It("backups and restores a cf", func() {
		if config.DeploymentToBackup.Name == config.DeploymentToRestore.Name {
			envsAreSame = true
		} else {
			printEnvsAreDifferentWarning()
		}

		// ### populate state in environment to be backed up
		for _, testCase := range testCases {
			testCase.BeforeBackup(config)
		}

		By("backing up " + config.DeploymentToBackup.Name)
		Eventually(RunCommandSuccessfully(fmt.Sprintf(
			"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup",
			session.WorkspaceDir,
			session.BinaryPath,
			config.BoshConfig.BoshURL,
			session.CertificatePath,
			config.BoshConfig.BoshClient,
			config.BoshConfig.BoshClientSecret,
			config.DeploymentToBackup.Name,
		))).Should(gexec.Exit(0))

		Eventually(StatusCode(config.DeploymentToBackup.ApiUrl)).Should(Equal(200))

		for _, testCase := range testCases {
			testCase.AfterBackup(config)
		}

		By("restoring to " + config.DeploymentToRestore.Name)
		Eventually(RunCommandSuccessfully(fmt.Sprintf(
			"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s restore --artifact-path $(ls %s | grep %s | head -n 1)",
			session.WorkspaceDir,
			session.BinaryPath,
			config.BoshConfig.BoshURL,
			session.CertificatePath,
			config.BoshConfig.BoshClient,
			config.BoshConfig.BoshClientSecret,
			config.DeploymentToRestore.Name,
			session.WorkspaceDir,
			config.DeploymentToBackup.Name,
		))).Should(gexec.Exit(0))

		// ### check state in restored environment
		for _, testCase := range testCases {
			testCase.AfterRestore(config)
		}
	})

	AfterEach(func() {
		By("cleaning up the artifact")
		Eventually(RunCommandSuccessfully(fmt.Sprintf("cd %s && rm -fr %s",
			session.WorkspaceDir,
			config.DeploymentToRestore.Name,
		))).Should(gexec.Exit(0))

		// ### clean up backup environment
		for _, testCase := range testCases {
			testCase.Cleanup(config)
		}

		session.Cleanup()
	})
}

func printEnvsAreDifferentWarning() {
	fmt.Println("     --------------------------------------------------------")
	fmt.Println("     NOTE: this suite is currently configured to back up from")
	fmt.Println("     one environment and restore to a difference one. Make   ")
	fmt.Println("     sure this is the intended configuration.                ")
	fmt.Println("     --------------------------------------------------------")
}
