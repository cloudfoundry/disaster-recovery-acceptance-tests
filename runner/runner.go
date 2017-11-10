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
	"io/ioutil"
	"path"
)

func RunDisasterRecoveryAcceptanceTests(configGetter ConfigGetter, testCases []TestCase) {
	var uniqueTestID string
	var testContext *TestContext
	var config Config
	var backupRunning bool
	var cfHomeTmpDir string
	var err error

	BeforeEach(func() {
		backupRunning = false
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

		cfHomeTmpDir, err = ioutil.TempDir("", "drats-cf-home")
		Expect(err).NotTo(HaveOccurred())

		for _, testCase := range testCases {
			err := os.Mkdir(cfHomeDir(cfHomeTmpDir,testCase), 0700)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("backups and restores a cf", func() {
		if config.DeploymentToBackup.Name != config.DeploymentToRestore.Name {
			printDeploymentsAreDifferentWarning()
		}

		// ### populate state in environment to be backed up
		for _, testCase := range testCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir,testCase))
			testCase.BeforeBackup(config)
		}

		backupRunning = true
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
		backupRunning = false

		Eventually(StatusCode(config.DeploymentToBackup.ApiUrl), 5*time.Minute).Should(Equal(200))

		for _, testCase := range testCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir,testCase))
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
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir,testCase))
			testCase.AfterRestore(config)
		}
	})

	AfterEach(func() {
		if backupRunning {
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
		}

		//TODO: Can we delete this?
		By("cleaning up the artifact")
		Eventually(RunCommandSuccessfully(fmt.Sprintf("cd %s && rm -fr %s",
			testContext.WorkspaceDir,
			config.DeploymentToRestore.Name,
		))).Should(gexec.Exit(0))

		// ### clean up backup environment
		for _, testCase := range testCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir,testCase))
			testCase.Cleanup(config)
		}

		testContext.Cleanup()

		for _, testCase := range testCases {
			cfHomeDir := testCase.Name() + "-cf-home"
			err := os.RemoveAll(cfHomeDir)
			Expect(err).NotTo(HaveOccurred())
		}
	})
}
func cfHomeDir(root string, testCase TestCase) string {
	return path.Join(root, testCase.Name()+"-cf-home")
}

func printDeploymentsAreDifferentWarning() {
	fmt.Println("     --------------------------------------------------------")
	fmt.Println("     NOTE: this suite is currently configured to back up from")
	fmt.Println("     one deployment and restore to a different one. Make     ")
	fmt.Println("     sure this is the intended configuration.                ")
	fmt.Println("     --------------------------------------------------------")
}
