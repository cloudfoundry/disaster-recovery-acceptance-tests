package runner

import (
	"fmt"

	"time"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunDisasterRecoveryAcceptanceTests(boshConfig BoshConfig, testCases []TestCase) {
	var envsAreSame bool
	var uniqueTestID string
	var jumpBoxSession *Session
	var config Config

	BeforeEach(func() {
		config = readConfigFromBOSHManifest(boshConfig)

		SetDefaultEventuallyTimeout(30 * time.Minute)
		uniqueTestID = RandomStringNumber()
		jumpBoxSession = NewSession(uniqueTestID, boshConfig)
	})

	It("backups and restores a cf", func() {
		if MustHaveEnv("DEPLOYMENT_TO_BACKUP") == MustHaveEnv("DEPLOYMENT_TO_RESTORE") {
			envsAreSame = true
		} else {
			printEnvsAreDifferentWarning()
		}

		By("finding credentials for the deployment to backup")
		urlForDeploymentToBackup, _, _ := FindCredentialsFor(DeploymentToBackup(), boshConfig)

		// ### populate state in environment to be backed up
		for _, testCase := range testCases {
			testCase.BeforeBackup(config)
		}

		By("backing up " + MustHaveEnv("DEPLOYMENT_TO_BACKUP"))
		Eventually(RunCommandSuccessfully(fmt.Sprintf(
			"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup",
			jumpBoxSession.WorkspaceDir,
			jumpBoxSession.BinaryPath,
			config.BoshConfig.BoshURL,
			jumpBoxSession.CertificatePath,
			config.BoshConfig.BoshClient,
			config.BoshConfig.BoshClientSecret,
			MustHaveEnv("DEPLOYMENT_TO_BACKUP"),
		))).Should(gexec.Exit(0))

		Eventually(StatusCode(urlForDeploymentToBackup)).Should(Equal(200))

		for _, testCase := range testCases {
			testCase.AfterBackup(config)
		}

		By("restoring to " + MustHaveEnv("DEPLOYMENT_TO_RESTORE"))
		Eventually(RunCommandSuccessfully(fmt.Sprintf(
			"cd %s && %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s restore --artifact-path $(ls %s | grep %s | head -n 1)",
			jumpBoxSession.WorkspaceDir,
			jumpBoxSession.BinaryPath,
			config.BoshConfig.BoshURL,
			jumpBoxSession.CertificatePath,
			config.BoshConfig.BoshClient,
			config.BoshConfig.BoshClientSecret,
			MustHaveEnv("DEPLOYMENT_TO_RESTORE"),
			jumpBoxSession.WorkspaceDir,
			MustHaveEnv("DEPLOYMENT_TO_BACKUP"),
		))).Should(gexec.Exit(0))

		// ### check state in restored environment
		for _, testCase := range testCases {
			testCase.AfterRestore(config)
		}
	})

	AfterEach(func() {
		By("cleaning up the artifact")
		Eventually(RunCommandSuccessfully(fmt.Sprintf("cd %s && rm -fr %s",
			jumpBoxSession.WorkspaceDir,
			MustHaveEnv("DEPLOYMENT_TO_RESTORE"),
		))).Should(gexec.Exit(0))

		// ### clean up backup environment
		for _, testCase := range testCases {
			testCase.Cleanup(config)
		}

		jumpBoxSession.Cleanup()
	})
}

func printEnvsAreDifferentWarning() {
	fmt.Println("     --------------------------------------------------------")
	fmt.Println("     NOTE: this suite is currently configured to back up from")
	fmt.Println("     one environment and restore to a difference one. Make   ")
	fmt.Println("     sure this is the intended configuration.                ")
	fmt.Println("     --------------------------------------------------------")
}

func readConfigFromBOSHManifest(boshConfig BoshConfig) Config {
	urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup(), boshConfig)
	urlForDeploymentToRestore, usernameForDeploymentToRestore, passwordForDeploymentToRestore := FindCredentialsFor(DeploymentToRestore(), boshConfig)

	return Config{
		DeploymentToBackup: CloudFoundryConfig{
			ApiUrl:        urlForDeploymentToBackup,
			AdminUsername: usernameForDeploymentToBackup,
			AdminPassword: passwordForDeploymentToBackup,
		},
		DeploymentToRestore: CloudFoundryConfig{
			ApiUrl:        urlForDeploymentToRestore,
			AdminUsername: usernameForDeploymentToRestore,
			AdminPassword: passwordForDeploymentToRestore,
		},
		BoshConfig: boshConfig,
	}
}
