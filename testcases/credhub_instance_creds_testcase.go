package testcases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	. "github.com/cloudfoundry/disaster-recovery-acceptance-tests/runner"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/mod/modfile"
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
	By("creating a test org and space")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully("cf auth", config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)

	By("setting up a test app")
	goVersion := fmt.Sprintf("go%s", tc.GetGoVersionFromGoModFile())

	RunCommandSuccessfully("cf push " + "--no-start " + tc.appName + " -p " + tc.testAppFixturePath + " -b go_buildpack" + " -f " + tc.testAppFixturePath + "/manifest.yml")
	RunCommandSuccessfully("cf set-env " + tc.appName + " GOVERSION " + goVersion + " > /dev/null")
	RunCommandSuccessfully("cf set-env " + tc.appName + " CREDHUB_CLIENT " + config.CloudFoundryConfig.CredHubClient + " > /dev/null")
	RunCommandSuccessfully("cf set-env " + tc.appName + " CREDHUB_SECRET " + config.CloudFoundryConfig.CredHubSecret + " > /dev/null")
	RunCommandSuccessfully("cf start " + tc.appName)

	By("verifying that the app can create and list a test secret")
	tc.appURL = GetAppURL(tc.appName)

	appCreateResponse := Get(tc.appURL + "/create")
	defer appCreateResponse.Body.Close()

	body, _ := io.ReadAll(appCreateResponse.Body)
	fmt.Println(string(body))
	Expect(appCreateResponse.StatusCode).To(Equal(http.StatusCreated))

	appListResponse := Get(tc.appURL + "/list")
	defer appListResponse.Body.Close()

	response, err := io.ReadAll(appListResponse.Body)
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(strings.NewReader(string(response))).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(1))
}

func (tc *CfCredhubSSITestCase) AfterBackup(config Config) {
	By("creating a second secret")
	appCreateResponse := Get(tc.appURL + "/create")
	defer appCreateResponse.Body.Close()

	Expect(appCreateResponse.StatusCode).To(Equal(http.StatusCreated))

	appListResponse := Get(tc.appURL + "/list")
	defer appListResponse.Body.Close()

	Expect(json.NewDecoder(appListResponse.Body).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(2))
}

func (tc *CfCredhubSSITestCase) GetGoVersionFromGoModFile() string {
	file_bytes, err := os.ReadFile(fmt.Sprintf("%s/go.mod", tc.testAppFixturePath))
	Expect(err).NotTo(HaveOccurred())
	f, err := modfile.Parse(fmt.Sprintf("%s/go.mod", tc.testAppFixturePath), file_bytes, nil)
	Expect(err).NotTo(HaveOccurred())
	return f.Go.Version

}
func (tc *CfCredhubSSITestCase) EnsureAfterSelectiveRestore(config Config) {
	By("repushing the test app")
	goVersion := fmt.Sprintf("go%s", tc.GetGoVersionFromGoModFile())

	RunCommandSuccessfully("cf push " + "--no-start " + tc.appName + " -p " + tc.testAppFixturePath + " -b go_buildpack" + " -f " + tc.testAppFixturePath + "/manifest.yml")
	RunCommandSuccessfully("cf set-env " + tc.appName + " GOVERSION " + goVersion + " > /dev/null")
	RunCommandSuccessfully("cf set-env " + tc.appName + " CREDHUB_CLIENT " + config.CloudFoundryConfig.CredHubClient + " > /dev/null")
	RunCommandSuccessfully("cf set-env " + tc.appName + " CREDHUB_SECRET " + config.CloudFoundryConfig.CredHubSecret + " > /dev/null")
	RunCommandSuccessfully("cf start " + tc.appName)
}

func (tc *CfCredhubSSITestCase) AfterRestore(config Config) {
	By("verifying that only one secret exists")
	appResponse := Get(tc.appURL + "/list")
	defer appResponse.Body.Close()

	Expect(json.NewDecoder(appResponse.Body).Decode(&listResponse)).To(Succeed())
	Expect(listResponse.Credentials).To(HaveLen(1))
}

func (tc *CfCredhubSSITestCase) Cleanup(config Config) {
	By("deleting the test secret")
	appResponse := Get(tc.appURL + "/clean")
	defer appResponse.Body.Close()

	By("deleting the test org")
	Expect(appResponse.StatusCode).To(Equal(http.StatusOK))
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}
