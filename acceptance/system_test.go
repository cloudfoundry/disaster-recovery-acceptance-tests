package acceptance

import (
	"fmt"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

var _ = Describe("backing up Cloud Foundry", func() {
	It("backups and restores a cf", func() {
		By("finding credentials for the deployment to backup")
		urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

		By("creating new orgs and spaces")
		RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
		RunCommandSuccessfully("cf create-org acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf create-space acceptance-test-space-" + uniqueTestID + " -o acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf target -s acceptance-test-space-" + uniqueTestID + " -o acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf push test_app_" + uniqueTestID + " -p " + testAppPath)

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

		Eventually(statusCode(urlForDeploymentToBackup)).Should(Equal(200))

		By("changing " + MustHaveEnv("DEPLOYMENT_TO_BACKUP"))
		RunCommandSuccessfully("cf delete -f test_app_" + uniqueTestID)
		RunCommandSuccessfully("cf delete-space -f acceptance-test-space-" + uniqueTestID)
		RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + uniqueTestID)

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

		By("finding credentials for the deployment to restore")
		urlForDeploymentToRestore, usernameForDeploymentToRestore, passwordForDeploymentToRestore := FindCredentialsFor(DeploymentToBackup())
		RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToRestore, "-u", usernameForDeploymentToRestore, "-p", passwordForDeploymentToRestore)

		By("verifying apps are back")
		RunCommandSuccessfully("cf target -s acceptance-test-space-" + uniqueTestID + " -o acceptance-test-org-" + uniqueTestID)
		url := GetAppUrl("test_app_" + uniqueTestID)

		Eventually(statusCode("https://"+url), 5*time.Minute, 5*time.Second).Should(Equal(200))

		By("verify orgs and spaces have been re-created")
		RunCommandSuccessfully("cf org acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf target -o acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf space acceptance-test-space-" + uniqueTestID)
	})

	AfterEach(func() {
		By("cleaning up the artifact")
		Eventually(RunOnJumpboxAsVcap(fmt.Sprintf("cd %s; rm -fr %s",
			jumpBoxSession.WorkspaceDir,
			MustHaveEnv("DEPLOYMENT_TO_RESTORE"),
		))).Should(gexec.Exit(0))

		By("cleaning up orgs and spaces")
		urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

		RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
		RunCommandSuccessfully("cf target -o acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf delete-space -f acceptance-test-space-" + uniqueTestID)
		RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + uniqueTestID)
	})
})
