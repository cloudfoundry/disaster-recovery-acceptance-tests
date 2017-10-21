package runner

import (
	"fmt"
	"strconv"

	"time"

	"os"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunDisasterRecoveryAcceptanceTests(configGetter ConfigGetter, testCases []TestCase) {
	var uniqueTestID string
	var testContext *TestContext
	var config Config

	BeforeEach(func() {
		config = configGetter.FindConfig()

		timeout := os.Getenv("DEFAULT_TIMEOUT_MINS")
		if timeout != "" {
			timeoutInt, err := strconv.Atoi(timeout)
			SetDefaultEventuallyTimeout(time.Duration(timeoutInt) * time.Minute)
			if err != nil {
				panic(fmt.Sprint("DEFAULT_TIMEOUT_MINS, if set, must be an integer\n"))
			}
		} else {
			SetDefaultEventuallyTimeout(15 * time.Minute)
		}

		uniqueTestID = RandomStringNumber()
		testContext = NewTestContext(uniqueTestID, config.BoshConfig)
	})

	It("backups and restores a cf", func() {
		if config.DeploymentToBackup.Name != config.DeploymentToRestore.Name {
			printDeploymentsAreDifferentWarning()
		}

		// ### populate state in environment to be backed up
		for _, testCase := range testCases {
			testCase.BeforeBackup(config)
		}

		By("backing up " + config.DeploymentToBackup.Name)
		Eventually(RunCommandSuccessfully(fmt.Sprintf(
			"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup",
			testContext.WorkspaceDir,
			testContext.BinaryPath,
			config.BoshConfig.BoshURL,
			testContext.CertificatePath,
			config.BoshConfig.BoshClient,
			config.BoshConfig.BoshClientSecret,
			config.DeploymentToBackup.Name,
		))).Should(gexec.Exit(0))

		Eventually(StatusCode(config.DeploymentToBackup.ApiUrl), 5*time.Minute).Should(Equal(200))

		for _, testCase := range testCases {
			testCase.AfterBackup(config)
		}

		By("restoring to " + config.DeploymentToRestore.Name)
		Eventually(RunCommandSuccessfully(fmt.Sprintf(
			"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s "+
				"--deployment %s restore --artifact-path $(ls %s | grep %s | head -n 1)",
			testContext.WorkspaceDir,
			testContext.BinaryPath,
			config.BoshConfig.BoshURL,
			testContext.CertificatePath,
			config.BoshConfig.BoshClient,
			config.BoshConfig.BoshClientSecret,
			config.DeploymentToRestore.Name,
			testContext.WorkspaceDir,
			config.DeploymentToBackup.Name,
		))).Should(gexec.Exit(0))

		// ### check state in restored environment
		for _, testCase := range testCases {
			testCase.AfterRestore(config)
		}
	})

	AfterEach(func() {
		By("running bbr backup-cleanup")
		Eventually(RunCommandSuccessfully(fmt.Sprintf(
			"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup-cleanup",
			testContext.WorkspaceDir,
			testContext.BinaryPath,
			config.BoshConfig.BoshURL,
			testContext.CertificatePath,
			config.BoshConfig.BoshClient,
			config.BoshConfig.BoshClientSecret,
			config.DeploymentToBackup.Name,
		))).Should(gexec.Exit(0))

		//TODO: Can we delete this?
		By("cleaning up the artifact")
		Eventually(RunCommandSuccessfully(fmt.Sprintf("cd %s && rm -fr %s",
			testContext.WorkspaceDir,
			config.DeploymentToRestore.Name,
		))).Should(gexec.Exit(0))

		// ### clean up backup environment
		for _, testCase := range testCases {
			testCase.Cleanup(config)
		}

		testContext.Cleanup()
	})
}

func printDeploymentsAreDifferentWarning() {
	fmt.Println("     --------------------------------------------------------")
	fmt.Println("     NOTE: this suite is currently configured to back up from")
	fmt.Println("     one deployment and restore to a different one. Make     ")
	fmt.Println("     sure this is the intended configuration.                ")
	fmt.Println("     --------------------------------------------------------")
}
