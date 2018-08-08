package testcases

import (
	"encoding/json"
	"net/http"
	"path"
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/gomega"
	"strings"
	"fmt"
	"io/ioutil"
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
	Credentials []struct{
		Name string
	}
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
	RunCommandSuccessfully("cf push " + tc.appName + " -p " + testAppPath + " -b go_buildpack" + " -f " + testAppPath + "/manifest.yml")

	tc.appURL = GetAppUrl(tc.appName)

	body := fmt.Sprintf(`{"credhub_client":"%s","credhub_secret": "%s"}`, config.CloudFoundryConfig.CredHubClient, config.CloudFoundryConfig.CredHubSecret)
	appResp := Post(tc.appURL+"/permission", "application/json", strings.NewReader(body))
	defer appResp.Body.Close()

	if appResp.StatusCode != http.StatusOK {
		resp, _ := ioutil.ReadAll(appResp.Body)
		Expect(appResp.StatusCode).To(Equal(http.StatusOK), string(resp))
	}

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
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
