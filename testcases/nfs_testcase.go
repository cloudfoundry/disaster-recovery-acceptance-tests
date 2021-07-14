package testcases

import (
	"fmt"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
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

func (tc *NFSTestCase) CheckDeployment(config Config) {
	By("checking if the NFS service is registered")
	RunCommandAndRetry(CF_CLI+" api --skip-ssl-validation", 3, config.CloudFoundryConfig.APIURL)
	RunCommandAndRetry(CF_CLI+" auth", 3, config.CloudFoundryConfig.AdminUsername, config.CloudFoundryConfig.AdminPassword)
	RunCommandSuccessfullyWithFailureMessage(
		tc.Name()+" test case cannot be run: NFS service is not registered",
		CF_CLI+" service-access -e "+config.CloudFoundryConfig.NFSServiceName,
	)
}

func (tc *NFSTestCase) BeforeBackup(config Config) {
	By("checking the service name and plane name are provided")
	Expect(config.CloudFoundryConfig.NFSServiceName).NotTo(BeEmpty(), "required config NFS service name not set")
	Expect(config.CloudFoundryConfig.NFSPlanName).NotTo(BeEmpty(), "required config NFS plan name not set")

	By("creating an NFS service broker and service instance")
	RunCommandSuccessfully(CF_CLI+" api --skip-ssl-validation", config.CloudFoundryConfig.APIURL)
	RunCommandSuccessfully(CF_CLI+" login --skip-ssl-validation -a", config.CloudFoundryConfig.APIURL,
		"-u", config.CloudFoundryConfig.AdminUsername, "-p", config.CloudFoundryConfig.AdminPassword)
	orgName := "acceptance-test-org-" + tc.uniqueTestID
	spaceName := "acceptance-test-space-" + tc.uniqueTestID
	RunCommandSuccessfully(CF_CLI + " create-org " + orgName)
	RunCommandSuccessfully(CF_CLI + " create-space " + spaceName + " -o " + orgName)
	RunCommandSuccessfully(CF_CLI + " target -o " + orgName + " -s " + spaceName)
	RunCommandSuccessfully(CF_CLI + " push dratsApp --docker-image docker/httpd --no-start --random-route")

	if config.CloudFoundryConfig.NFSCreateServiceBroker {
		RunCommandSuccessfully(CF_CLI + " create-service-broker nfsbroker-drats-" + tc.uniqueTestID + " " +
			config.CloudFoundryConfig.NFSBrokerUser + " " + config.CloudFoundryConfig.NFSBrokerPassword + " " +
			config.CloudFoundryConfig.NFSBrokerURL)
	}

	RunCommandSuccessfully(CF_CLI + " enable-service-access " + config.CloudFoundryConfig.NFSServiceName + " -o " + orgName)
	RunCommandSuccessfully(CF_CLI + " create-service " + config.CloudFoundryConfig.NFSServiceName + " " +
		config.CloudFoundryConfig.NFSPlanName + " " + tc.instanceName + " -c " +
		`'{"share":"someserver.someplace.com/someshare"}'`)
}

func (tc *NFSTestCase) AfterBackup(config Config) {
	By("deleting the NFS service instance after backup")
	RunCommandSuccessfully(CF_CLI + " delete-service " + tc.instanceName + " -f")
}

func (tc *NFSTestCase) EnsureAfterSelectiveRestore(config Config) {
	By("repushing apps if restoring from a selective restore")
	RunCommandSuccessfully(CF_CLI + " push dratsApp --docker-image docker/httpd --no-start --random-route")
}

func (tc *NFSTestCase) AfterRestore(config Config) {
	By("re-binding the NFS service instance after restore")
	RunCommandSuccessfully(CF_CLI + " bind-service dratsApp " + tc.instanceName + ` -c '{"uid":5000,"gid":5000}'`)
}

func (tc *NFSTestCase) Cleanup(config Config) {
	By("nfs cleanup")
	RunCommandSuccessfully(CF_CLI + " delete-org -f acceptance-test-org-" + tc.uniqueTestID)

	if config.CloudFoundryConfig.NFSCreateServiceBroker {
		RunCommandSuccessfully(CF_CLI + " delete-service-broker -f nfsbroker-drats-" + tc.uniqueTestID)
	}
}
