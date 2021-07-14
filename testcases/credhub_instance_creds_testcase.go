package testcases

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/gomega"
)

type CfCredhubSSITestCase struct {
	uniqueTestID       string
	name               string
	appName            string
	appURL             string
	testAppFixturePath string
}

func NewCfCredhubSSITestCase() *CfCredhubSSITestCase {
	id := RandomStringNumber()

	credhubAppPath, appPathPresent := os.LookupEnv("CREDHUB_APP_PATH")

	if no(appPathPresent) {
		credhubAppPath = path.Join(CurrentTestDir(), "/../fixtures/credhub-test-app")
	}

	return &CfCredhubSSITestCase{
		uniqueTestID:       id,
		name:               "cf-credhub",
		appName:            "app" + id,
		testAppFixturePath: credhubAppPath,
	}
}

func no(b bool) bool {
	return !b
}

var listResponse struct {
	Credentials []struct {
		Name string
	}
}

func (tc *CfCredhubSSITestCase) Name() string {
	return tc.name
}

func (tc *CfCredhubSSITestCase) CheckDeployment(config Config) {
}

func (tc *CfCredhubSSITestCase) BeforeBackup(config Config) {
	RunCommandSuccessfully(CF_CLI+" api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully(CF_CLI+" auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully(CF_CLI + " create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully(CF_CLI + " target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)

	RunCommandSuccessfully(CF_CLI + " push " + "--no-start " + tc.appName + " -p " + tc.testAppFixturePath + " -b go_buildpack" + " -f " + tc.testAppFixturePath + "/manifest.yml")
	RunCommandSuccessfully(CF_CLI + " set-env " + tc.appName + " CREDHUB_CLIENT " + config.CloudFoundryConfig.CredHubClient + " > /dev/null")
	RunCommandSuccessfully(CF_CLI + " set-env " + tc.appName + " CREDHUB_SECRET " + config.CloudFoundryConfig.CredHubSecret + " > /dev/null")
	RunCommandSuccessfully(CF_CLI + " start " + tc.appName)

	tc.appURL = GetAppURL(tc.appName)

	appCreateResponse := Get(tc.appURL + "/create")
	defer appCreateResponse.Body.Close()

	body, _ := ioutil.ReadAll(appCreateResponse.Body)
	fmt.Println(string(body))
	Expect(appCreateResponse.StatusCode).To(Equal(http.StatusCreated))

	appListResponse := Get(tc.appURL + "/list")
	defer appListResponse.Body.Close()

	response, err := ioutil.ReadAll(appListResponse.Body)
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(strings.NewReader(string(response))).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(1))
}

func (tc *CfCredhubSSITestCase) AfterBackup(config Config) {
	appCreateResponse := Get(tc.appURL + "/create")
	defer appCreateResponse.Body.Close()

	Expect(appCreateResponse.StatusCode).To(Equal(http.StatusCreated))

	appListResponse := Get(tc.appURL + "/list")
	defer appListResponse.Body.Close()

	Expect(json.NewDecoder(appListResponse.Body).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(2))
}

func (tc *CfCredhubSSITestCase) EnsureAfterSelectiveRestore(config Config) {
	RunCommandSuccessfully(CF_CLI + " push " + "--no-start " + tc.appName + " -p " + tc.testAppFixturePath + " -b go_buildpack" + " -f " + tc.testAppFixturePath + "/manifest.yml")
	RunCommandSuccessfully(CF_CLI + " set-env " + tc.appName + " CREDHUB_CLIENT " + config.CloudFoundryConfig.CredHubClient + " > /dev/null")
	RunCommandSuccessfully(CF_CLI + " set-env " + tc.appName + " CREDHUB_SECRET " + config.CloudFoundryConfig.CredHubSecret + " > /dev/null")
	RunCommandSuccessfully(CF_CLI + " start " + tc.appName)
}

func (tc *CfCredhubSSITestCase) AfterRestore(config Config) {
	appResponse := Get(tc.appURL + "/list")
	defer appResponse.Body.Close()

	Expect(json.NewDecoder(appResponse.Body).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(1))
}

func (tc *CfCredhubSSITestCase) Cleanup(config Config) {
	appResponse := Get(tc.appURL + "/clean")
	defer appResponse.Body.Close()

	Expect(appResponse.StatusCode).To(Equal(http.StatusOK))
	RunCommandSuccessfully(CF_CLI + " delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
