package testcases

import (

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	"os"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type NFSTestCase struct {
	uniqueTestID string
	instanceName string
}

func NewNFSTestCases() *NFSTestCase {
	id := RandomStringNumber()
	name := fmt.Sprintf("service-instance-%s", id)
	return &NFSTestCase{uniqueTestID: id, instanceName: name}
}

func (tc *NFSTestCase) Name() string {
	return "cf-nfsbroker"
}

func (tc *NFSTestCase) BeforeBackup(config Config) {
	By("nfs before backup")
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	fmt.Printf("-------->cwd %s", cwd)
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.DeploymentToBackup.ApiUrl, "-u", config.DeploymentToBackup.AdminUsername, "-p", config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf push dratsApp --docker-image docker/httpd --no-start")

	RunCommandSuccessfully("cf create-service-broker " + "nfsbroker-drats-" + tc.uniqueTestID + " " + os.Getenv("BROKER_USER") + " " + os.Getenv("BROKER_PASSWORD") + " " + os.Getenv("BROKER_URL"))
	RunCommandSuccessfully("cf enable-service-access " + os.Getenv("SERVICE_NAME"))
	RunCommandSuccessfully("cf create-service " + os.Getenv("SERVICE_NAME")+ " " + os.Getenv("PLAN_NAME") + " " + tc.instanceName + " -c " + fmt.Sprintf("'{\"share\":\"%s%s\"}'", os.Getenv("SERVER_ADDRESS"), os.Getenv("SHARE")))
}

func (tc *NFSTestCase) AfterBackup(config Config) {
	By("nfs after backup")
	RunCommandSuccessfully("cf delete-service " + tc.instanceName + " -f")
}

func (tc *NFSTestCase) AfterRestore(config Config) {
	By("nfs after backup")
	RunCommandSuccessfully("cf bind-service dratsApp " + tc.instanceName + ` -c '{"uid":5000,"gid":5000}'`)
}

func (tc *NFSTestCase) Cleanup(config Config) {
	By("nfs cleanup")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf delete-service-broker -f " + "nfsbroker-drats-" + tc.uniqueTestID)
}

func (tc *NFSTestCase) deletePushedApps(config Config) {
}
