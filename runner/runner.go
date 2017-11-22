package runner

import (
	"fmt"
	"strconv"

	"time"

	"os"

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
	var filteredTestCases []TestCase

	BeforeEach(func() {
		focusedSuiteName := os.Getenv("FOCUSED_SUITE_NAME")
		skipSuiteName := os.Getenv("SKIP_SUITE_NAME")
		filteredTestCases = FilterTestCasesWithRegexes(testCases, skipSuiteName, focusedSuiteName)
		fmt.Println("Running testcases:")
		for _, testCase := range filteredTestCases {
			fmt.Println(testCase.Name())
		}

		backupRunning = false
		config = configGetter.FindConfig(filteredTestCases)

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

		for _, testCase := range filteredTestCases {
			err := os.Mkdir(cfHomeDir(cfHomeTmpDir, testCase), 0700)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("backups and restores a cf", func() {
		if config.DeploymentToBackup.Name != config.DeploymentToRestore.Name {
			printDeploymentsAreDifferentWarning()
		}

		// ### populate state in environment to be backed up
		for _, testCase := range filteredTestCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
			testCase.BeforeBackup(config)
		}

		backupRunning = true
		By("backing up " + config.DeploymentToBackup.Name)
		RunCommandSuccessfullyWithFailureMessage(
			"bbr deployment backup",
			fmt.Sprintf(
				"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup",
				testContext.WorkspaceDir,
				testContext.BinaryPath,
				config.BoshConfig.BoshURL,
				testContext.CertificatePath,
				config.BoshConfig.BoshClient,
				config.BoshConfig.BoshClientSecret,
				config.DeploymentToBackup.Name,
			))
		backupRunning = false

		Eventually(StatusCode(config.DeploymentToBackup.ApiUrl), 5*time.Minute).Should(Equal(200))

		for _, testCase := range filteredTestCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
			testCase.AfterBackup(config)
		}

		By("restoring to " + config.DeploymentToRestore.Name)
		RunCommandSuccessfullyWithFailureMessage(
			"bbr deployment restore",
			fmt.Sprintf(
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
			))

		// ### check state in restored environment
		for _, testCase := range filteredTestCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
			testCase.AfterRestore(config)
		}
	})

	AfterEach(func() {
		var backupCleanupSession, artifactCleanupSession *gexec.Session
		var removeHomeDirErr error
		if backupRunning {
			By("running bbr backup-cleanup")
			backupCleanupSession = RunCommandWithFailureMessage(
				"bbr deployment backup-cleanup",
				fmt.Sprintf(
					"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup-cleanup",
					testContext.WorkspaceDir,
					testContext.BinaryPath,
					config.BoshConfig.BoshURL,
					testContext.CertificatePath,
					config.BoshConfig.BoshClient,
					config.BoshConfig.BoshClientSecret,
					config.DeploymentToBackup.Name,
				))
		}
		//TODO: Can we delete this?
		By("cleaning up the artifact")
		artifactCleanupSession = RunCommand(fmt.Sprintf("cd %s && rm -fr %s",
			testContext.WorkspaceDir,
			config.DeploymentToRestore.Name,
		))

		By("running the individual test-case cleanup commands")
		for _, testCase := range filteredTestCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
			testCase.Cleanup(config)
		}

		By("removing individual test-case cf-home directory")
		for _, testCase := range filteredTestCases {
			removeHomeDirErr = os.RemoveAll(cfHomeDir(cfHomeTmpDir, testCase))
		}

		By("cleaning up the test context")
		testContext.Cleanup()

		if backupRunning {
			Expect(backupCleanupSession).To(gexec.Exit(0))
		}

		Expect(artifactCleanupSession).To(gexec.Exit(0))
		Expect(removeHomeDirErr).ToNot(HaveOccurred())

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
