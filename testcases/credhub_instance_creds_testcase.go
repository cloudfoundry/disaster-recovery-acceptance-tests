package testcases

import (
	"encoding/json"
	"net/http"
	"path"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/gomega"
)

type CfCredhubSSITestCase struct {
	uniqueTestID string
	name         string
	appName      string
	appURL       string
}

func NewCfCredhubSSITestCase() *CfCredhubSSITestCase {
	id := RandomStringNumber()
	return &CfCredhubSSITestCase{
		uniqueTestID: id,
		name:         "cf-credhub",
		appName:      "app" + id,
	}
}

var listResponse struct {
	Credentials []struct{}
}

func (tc *CfCredhubSSITestCase) Name() string {
	return tc.name
}

func (tc *CfCredhubSSITestCase) BeforeBackup(config Config) {
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.ApiUrl)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)

	var testAppPath = path.Join(CurrentTestDir(), "/../fixtures/credhub-test-app")
	RunCommandSuccessfully("cf push " + tc.appName + " -p " + testAppPath + " -f " + testAppPath + "/manifest.yml")

	tc.appURL = GetAppUrl(tc.appName)
	appResponse := Get(tc.appURL + "/create")
	Expect(appResponse.StatusCode).To(Equal(http.StatusCreated))

	appResponse = Get(tc.appURL + "/list")
	Expect(json.NewDecoder(appResponse.Body).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(1))
}

func (tc *CfCredhubSSITestCase) AfterBackup(config Config) {
	appResponse := Get(tc.appURL + "/create")
	Expect(appResponse.StatusCode).To(Equal(http.StatusCreated))

	appResponse = Get(tc.appURL + "/list")
	Expect(json.NewDecoder(appResponse.Body).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(2))
}

func (tc *CfCredhubSSITestCase) AfterRestore(config Config) {
	appResponse := Get(tc.appURL + "/list")
	Expect(json.NewDecoder(appResponse.Body).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(1))
}

func (tc *CfCredhubSSITestCase) Cleanup(config Config) {
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.ApiUrl)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
