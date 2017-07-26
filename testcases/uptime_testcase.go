package testcases

import (
	"net/http"
	"path"
	"time"

	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

type AppUptimeTestCase struct {
	uniqueTestID            string
	stopCheckingAppAlive    chan<- bool
	stopCheckingAPIGoesDown chan<- bool
	valueApiWasDown         <-chan bool
}

func NewAppUptimeTestCase() *AppUptimeTestCase {
	id := RandomStringNumber()
	return &AppUptimeTestCase{uniqueTestID: id}
}

func (tc *AppUptimeTestCase) BeforeBackup() {
	urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

	RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID + " -s acceptance-test-space-" + tc.uniqueTestID)
	var testAppFixturePath = path.Join(CurrentTestDir(), "/../fixtures/test_app/")
	RunCommandSuccessfully("cf push test_app_" + tc.uniqueTestID + " -p " + testAppFixturePath)

	By("checking the app stays up")
	appUrl := GetAppUrl("test_app_" + tc.uniqueTestID)
	tc.stopCheckingAppAlive = checkAppRemainsAlive(appUrl)
	tc.stopCheckingAPIGoesDown, tc.valueApiWasDown = checkApiGoesDown()
}

func (tc *AppUptimeTestCase) AfterBackup() {
	By("stopping checking the app")
	log.Println("writing to stopCheckingAppAlive...")
	tc.stopCheckingAppAlive <- true
	log.Println("writing to stopCheckingAPIGoesDown...")
	tc.stopCheckingAPIGoesDown <- true
	log.Println("reading from valueApiWasDown...")
	Expect(<-tc.valueApiWasDown).To(BeTrue())
}

func (tc *AppUptimeTestCase) AfterRestore() {

}

func (tc *AppUptimeTestCase) Cleanup() {
	By("cleaning up orgs, spaces and apps")
	urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

	RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf delete-space -f acceptance-test-space-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}

func checkApiGoesDown() (chan<- bool, <-chan bool) {
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

func checkAppRemainsAlive(url string) chan<- bool {
	doneChannel := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)
	tickerChannel := ticker.C

	go func() {
		defer GinkgoRecover()
		for {
			select {
			case <-doneChannel:
				ticker.Stop()
				return
			case <-tickerChannel:
				Expect(Get(url).StatusCode).To(Equal(http.StatusOK))
			}
		}
	}()

	return doneChannel
}
