package acceptance

import (
	"net/http"

	"time"

	"fmt"

	"crypto/tls"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

var _ = Describe("backing up Cloud Foundry", func() {
	It("apps remain reachable for the duration", func() {
		By("finding credentials for the deployment to backup")
		urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

		By("creating new orgs, spaces and apps")
		RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
		RunCommandSuccessfully("cf create-org acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf create-space acceptance-test-space-" + uniqueTestID + " -o acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf target -o acceptance-test-org-" + uniqueTestID + " -s acceptance-test-space-" + uniqueTestID)
		RunCommandSuccessfully("cf push test_app_" + uniqueTestID + " -p " + testAppPath)

		By("checking the app stays up")
		appUrl := GetAppUrl("test_app_" + uniqueTestID)
		stopCheckingAppAlive := CheckAppRemainsAlive(appUrl)
		stopCheckingAPIGoesDown, apiWasDown := CheckApiGoesDown()

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

		By("stopping checking the app")
		stopCheckingAppAlive <- true
		stopCheckingAPIGoesDown <- true
		Expect(<-apiWasDown).To(BeTrue())
	})

	AfterEach(func() {
		By("cleaning up the artifact")
		Eventually(RunOnJumpboxAsVcap(fmt.Sprintf("cd %s; rm -fr %s",
			jumpBoxSession.WorkspaceDir,
			MustHaveEnv("DEPLOYMENT_TO_BACKUP"),
		))).Should(gexec.Exit(0))

		By("cleaning up orgs, spaces and apps")
		urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

		RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
		RunCommandSuccessfully("cf target -o acceptance-test-org-" + uniqueTestID)
		RunCommandSuccessfully("cf delete-space -f acceptance-test-space-" + uniqueTestID)
		RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + uniqueTestID)
	})
})

func statusCode(url string) func() (int, error) {
	return func() (int, error) {
		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
		resp, err := client.Get(url)
		if err != nil {
			return 0, err
		}
		return resp.StatusCode, err
	}
}

func CheckApiGoesDown() (chan<- bool, <-chan bool) {
	doneChannel := make(chan bool)
	valueApiWasDown := make(chan bool)
	ticker := time.NewTicker(1 * time.Second)
	tickerChannel := ticker.C
	go func() {
		var apiWasDown bool
		defer GinkgoRecover()
		for {
			select {
			case <-doneChannel:
				valueApiWasDown <- apiWasDown
				return
			case <-tickerChannel:
				if RunCommand("cf orgs").ExitCode() == 1 {
					apiWasDown = true
					ticker.Stop()
				}
			}
		}
	}()
	return doneChannel, valueApiWasDown
}

func CheckAppRemainsAlive(url string) chan<- bool {
	doneChannel := make(chan bool, 1)
	tickerChannel := time.NewTicker(1 * time.Second).C

	go func() {
		defer GinkgoRecover()
		for {
			select {
			case <-doneChannel:
				return
			case <-tickerChannel:
				Expect(Get(url).StatusCode).To(Equal(http.StatusOK))
			}
		}
	}()

	return doneChannel
}
