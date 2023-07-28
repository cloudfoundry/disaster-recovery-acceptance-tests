package runner

import (
	"fmt"
	"time"

	"os"

	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunDisasterRecoveryAcceptanceTests(config Config, testCases []TestCase) {
	var uniqueTestID string
	var testContext *TestContext
	var backupRunning, restoreRunning bool
	var cfHomeTmpDir string
	var err error

	BeforeEach(func() {
		SetDefaultEventuallyTimeout(config.Timeout)

		fmt.Println("\nCF CLI version:")
		RunCommandSuccessfully("cf", "version")

		fmt.Println("\nRunning test cases:")
		for _, testCase := range testCases {
			fmt.Println(testCase.Name())
		}

		fmt.Println("\nChecking deployment has been set up for test cases...")
		for _, testCase := range testCases {
			testCase.CheckDeployment(config)
		}

		backupRunning = false

		uniqueTestID = RandomStringNumber()
		testContext, err = NewTestContext(uniqueTestID, config.BoshConfig)
		Expect(err).NotTo(HaveOccurred())

		cfHomeTmpDir = GinkgoT().TempDir()
		for _, testCase := range testCases {
			err := os.Mkdir(cfHomeDir(cfHomeTmpDir, testCase), 0700)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("backs up and restores a CF", func() {
		By("populating state in environment to be backed up")
		for _, testCase := range testCases {
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
			By("running the BeforeBackup step for " + testCase.Name())
			testCase.BeforeBackup(config)
		}

		backupRunning = true
		By("backing up " + config.CloudFoundryConfig.Name)
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
				config.CloudFoundryConfig.Name,
			))
		backupRunning = false

		Eventually(StatusCode(config.CloudFoundryConfig.APIURL), 5*time.Minute).Should(Equal(200))

		for _, testCase := range testCases {
			By("running the AfterBackup step for " + testCase.Name())
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
			testCase.AfterBackup(config)
		}

		if config.DeleteAndRedeployCF {
			By("deleting the deployment")
			manifestSession := RunCommandSuccessfullySilently("bosh-cli",
				"-e", config.BoshURL,
				"--ca-cert", testContext.CertificatePath,
				"--client", config.BoshConfig.BoshClient,
				"--client-secret", config.BoshConfig.BoshClientSecret,
				"-d", config.CloudFoundryConfig.Name,
				"manifest",
			)
			file, err := os.CreateTemp("", "cf")
			Expect(err).NotTo(HaveOccurred())
			_, err = file.Write(manifestSession.Out.Contents())
			Expect(err).NotTo(HaveOccurred())

			RunCommandSuccessfully("bosh-cli", "-n",
				"-e", config.BoshURL,
				"--ca-cert", testContext.CertificatePath,
				"--client", config.BoshConfig.BoshClient,
				"--client-secret", config.BoshConfig.BoshClientSecret,
				"-d", config.CloudFoundryConfig.Name,
				"delete-deployment",
			)

			RunCommandSuccessfully("bosh-cli", "-n",
				"-e", config.BoshURL,
				"--ca-cert", testContext.CertificatePath,
				"--client", config.BoshConfig.BoshClient,
				"--client-secret", config.BoshConfig.BoshClientSecret,
				"-d", config.CloudFoundryConfig.Name,
				"deploy", file.Name(),
			)
		}

		restoreRunning = true
		By("restoring to " + config.CloudFoundryConfig.Name)
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
				config.CloudFoundryConfig.Name,
				testContext.WorkspaceDir,
				config.CloudFoundryConfig.Name,
			))
		restoreRunning = false

		Eventually(StatusCode(config.CloudFoundryConfig.APIURL), 5*time.Minute).Should(Equal(200))

		if config.SelectiveBackup {
			for _, testCase := range testCases {
				By("running the EnsureAfterSelectiveRestore step for " + testCase.Name())
				os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
				testCase.EnsureAfterSelectiveRestore(config)
			}
		}

		By("checking state in restored environment")
		for _, testCase := range testCases {
			By("running the AfterRestore step for " + testCase.Name())
			os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
			testCase.AfterRestore(config)
		}
	})

	AfterEach(func() {
		if backupRunning {
			By("running bbr backup-cleanup")
			backupCleanupSession := RunCommandWithFailureMessage(
				"bbr deployment backup-cleanup",
				fmt.Sprintf(
					"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup-cleanup",
					testContext.WorkspaceDir,
					testContext.BinaryPath,
					config.BoshConfig.BoshURL,
					testContext.CertificatePath,
					config.BoshConfig.BoshClient,
					config.BoshConfig.BoshClientSecret,
					config.CloudFoundryConfig.Name,
				))
			Expect(backupCleanupSession).To(gexec.Exit(0))
		}
	})

	AfterEach(func() {
		if restoreRunning {
			By("running bbr restore-cleanup")
			restoreCleanupSession := RunCommandWithFailureMessage(
				"bbr deployment restore-cleanup",
				fmt.Sprintf(
					"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s restore-cleanup",
					testContext.WorkspaceDir,
					testContext.BinaryPath,
					config.BoshConfig.BoshURL,
					testContext.CertificatePath,
					config.BoshConfig.BoshClient,
					config.BoshConfig.BoshClientSecret,
					config.CloudFoundryConfig.Name,
				))
			Expect(restoreCleanupSession).To(gexec.Exit(0))
		}
	})

	AfterEach(func() {
		By("cleaning up the artifact")
		artifactCleanupSession := RunCommand(fmt.Sprintf("cd %s && rm -fr %s",
			testContext.WorkspaceDir,
			config.CloudFoundryConfig.Name,
		))
		Expect(artifactCleanupSession).To(gexec.Exit(0))
	})

	for _, testCase := range testCases {
		func(testCase TestCase) {
			AfterEach(func() {
				By("running the Cleanup step for " + testCase.Name())
				os.Setenv("CF_HOME", cfHomeDir(cfHomeTmpDir, testCase))
				testCase.Cleanup(config)
			})
		}(testCase)
	}

	AfterEach(func() {
		By("cleaning up the test context")
		testContext.Cleanup()
	})
}

func cfHomeDir(root string, testCase TestCase) string {
	return path.Join(root, testCase.Name()+"-cf-home")
}
