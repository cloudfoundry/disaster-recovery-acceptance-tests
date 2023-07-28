package testcases

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	. "github.com/cloudfoundry/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const CURL_ERROR_FOR_404 = 22

type AppUptimeTestCase struct {
	uniqueTestID            string
	stopCheckingAppAlive    chan<- bool
	stopCheckingAPIGoesDown chan<- bool
	valueAPIWasDown         <-chan bool
	name                    string
}

func NewAppUptimeTestCase() *AppUptimeTestCase {
	id := RandomStringNumber()
	return &AppUptimeTestCase{uniqueTestID: id, name: "app-uptime"}
}

func (tc *AppUptimeTestCase) Name() string {
	return tc.name
}

func (tc *AppUptimeTestCase) CheckDeployment(config Config) {
}

func (tc *AppUptimeTestCase) BeforeBackup(config Config) {
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID + " -s acceptance-test-space-" + tc.uniqueTestID)
	var testAppFixturePath = path.Join(CurrentTestDir(), "/../fixtures/test_app/")
	RunCommandSuccessfully("cf push test_app_" + tc.uniqueTestID + " -p " + testAppFixturePath)

	By("checking the app stays up")
	appURL := GetAppURL("test_app_" + tc.uniqueTestID)
	tc.stopCheckingAppAlive = checkAppRemainsAlive(appURL)
	tc.stopCheckingAPIGoesDown, tc.valueAPIWasDown = checkAPIGoesDown(config.CloudFoundryConfig.APIURL)
}

func (tc *AppUptimeTestCase) AfterBackup(config Config) {
	By("stopping checking the app")
	log.Println("writing to stopCheckingAppAlive...")
	tc.stopCheckingAppAlive <- true
	log.Println("writing to stopCheckingAPIGoesDown...")
	tc.stopCheckingAPIGoesDown <- true
	log.Println("reading from valueAPIWasDown...")
	Expect(<-tc.valueAPIWasDown).To(BeTrue(), "expected api to be down, but it isn't")
}

func (tc *AppUptimeTestCase) EnsureAfterSelectiveRestore(config Config) {}

func (tc *AppUptimeTestCase) AfterRestore(config Config) {

}

func (tc *AppUptimeTestCase) Cleanup(config Config) {
	By("cleaning up orgs, spaces and apps")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}

func checkAPIGoesDown(apiURL string) (chan<- bool, <-chan bool) {
	doneChannel := make(chan bool)
	valueAPIWasDown := make(chan bool)
	ticker := time.NewTicker(1 * time.Second)
	tickerChannel := ticker.C

	go func() {
		var apiWasDown bool
		defer GinkgoRecover()
		for {
			select {
			case <-doneChannel:
				valueAPIWasDown <- apiWasDown
				return
			case <-tickerChannel:
				if RunCommand("curl", "-k", "--fail", "--max-time", "1", apiURL, " 2>/dev/null > /dev/null").ExitCode() == CURL_ERROR_FOR_404 {
					apiWasDown = true
					ticker.Stop()
				}
			}
		}
	}()

	return doneChannel, valueAPIWasDown
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
				resp := Get(url)
				Expect(resp.StatusCode).To(Equal(http.StatusOK), fmt.Sprintf("%s - expected app to consistently respond 200 OK during backup", time.Now().UTC()))
				resp.Body.Close()
			}
		}
	}()

	return doneChannel
}
