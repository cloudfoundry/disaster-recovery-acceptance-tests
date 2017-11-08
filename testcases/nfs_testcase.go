package testcases

import (
	"fmt"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	. "github.com/onsi/ginkgo"
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
	By("creating an NFS service broker and service instance")
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.DeploymentToBackup.ApiUrl)
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.DeploymentToBackup.ApiUrl,
		"-u", config.DeploymentToBackup.AdminUsername, "-p", config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID +
		" -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID +
		" -o acceptance-test-org-" + tc.uniqueTestID)
	routePath := fmt.Sprintf("dratsapp-%s.%s", tc.uniqueTestID, config.DeploymentToBackup.ApiUrl)
	RunCommandSuccessfully("cf push dratsApp --docker-image docker/httpd --no-start --route-path " + routePath)

	if config.DeploymentToBackup.NFSBrokerUser != "" {
		RunCommandSuccessfully("cf create-service-broker " + "nfsbroker-drats-" + tc.uniqueTestID + " " +
			config.DeploymentToBackup.NFSBrokerUser + " " + config.DeploymentToBackup.NFSBrokerPassword + " " +
			config.DeploymentToBackup.NFSBrokerUrl)
	}
	RunCommandSuccessfully("cf enable-service-access " + config.DeploymentToBackup.NFSServiceName)
	RunCommandSuccessfully("cf create-service " + config.DeploymentToBackup.NFSServiceName + " " +
		config.DeploymentToBackup.NFSPlanName + " " + tc.instanceName + " -c " +
		`'{"share":"someserver.someplace.com/someshare"}'`)
}

func (tc *NFSTestCase) AfterBackup(config Config) {
	By("deleting the NFS service instance after backup")
	RunCommandSuccessfully("cf delete-service " + tc.instanceName + " -f")
}

func (tc *NFSTestCase) AfterRestore(config Config) {
	By("re-binding the NFS service instance after restore")
	RunCommandSuccessfully("cf bind-service dratsApp " + tc.instanceName + ` -c '{"uid":5000,"gid":5000}'`)
}

func (tc *NFSTestCase) Cleanup(config Config) {
	By("nfs cleanup")
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
	if config.DeploymentToBackup.NFSBrokerUser != "" {
		RunCommandSuccessfully("cf delete-service-broker -f " + "nfsbroker-drats-" + tc.uniqueTestID)
	}
}

func (tc *NFSTestCase) deletePushedApps(config Config) {
}
