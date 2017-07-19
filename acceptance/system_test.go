package acceptance

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

var _ = Describe("backing up Cloud Foundry", func() {
	var envsAreSame bool

	It("backups and restores a cf", func() {
		if MustHaveEnv("DEPLOYMENT_TO_BACKUP") == MustHaveEnv("DEPLOYMENT_TO_RESTORE") {
			envsAreSame = true
		} else {
			printEnvsAreDifferentWarning()
		}

		By("finding credentials for the deployment to backup")
		urlForDeploymentToBackup, _, _ := FindCredentialsFor(DeploymentToBackup())

		// ### populate state in environment to be backed up
		for _, testCase := range testCases {
			testCase.PopulateState()
		}

		By("backing up " + MustHaveEnv("DEPLOYMENT_TO_BACKUP"))
		Eventually(RunOnJumpboxAsVcap(fmt.Sprintf(
			"cd %s; %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s backup",
			jumpBoxSession.WorkspaceDir,
			jumpBoxSession.BinaryPath,
			MustHaveEnv("BOSH_URL"),
			jumpBoxSession.CertificatePath,
			MustHaveEnv("BOSH_CLIENT"),
			MustHaveEnv("BOSH_CLIENT_SECRET"),
			MustHaveEnv("DEPLOYMENT_TO_BACKUP"),
		))).Should(gexec.Exit(0))

		Eventually(StatusCode(urlForDeploymentToBackup)).Should(Equal(200))

		if envsAreSame {
			// ### clean up state in backed up environment
			for _, testCase := range testCases {
				testCase.Cleanup()
			}
		}

		By("restoring to " + MustHaveEnv("DEPLOYMENT_TO_RESTORE"))
		Eventually(RunOnJumpboxAsVcap(fmt.Sprintf(
			"cd %s; %s deployment --target %s --ca-cert %s --username %s --password %s --deployment %s restore --artifact-path $(ls %s | grep %s | head -n 1)",
			jumpBoxSession.WorkspaceDir,
			jumpBoxSession.BinaryPath,
			MustHaveEnv("BOSH_URL"),
			jumpBoxSession.CertificatePath,
			MustHaveEnv("BOSH_CLIENT"),
			MustHaveEnv("BOSH_CLIENT_SECRET"),
			MustHaveEnv("DEPLOYMENT_TO_RESTORE"),
			jumpBoxSession.WorkspaceDir,
			MustHaveEnv("DEPLOYMENT_TO_BACKUP"),
		))).Should(gexec.Exit(0))

		// ### check state in restored environment
		for _, testCase := range testCases {
			testCase.CheckState()
		}
	})

	AfterEach(func() {
		By("cleaning up the artifact")
		Eventually(RunOnJumpboxAsVcap(fmt.Sprintf("cd %s; rm -fr %s",
			jumpBoxSession.WorkspaceDir,
			MustHaveEnv("DEPLOYMENT_TO_RESTORE"),
		))).Should(gexec.Exit(0))

		// ### clean up backup environment
		for _, testCase := range testCases {
			testCase.Cleanup()
		}
	})
})
