package disaster_recovery_acceptance_tests

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("backing up Cloud Foundry", func() {
	focusedSuiteName := os.Getenv("FOCUSED_SUITE_NAME")
	skipSuiteName := os.Getenv("SKIP_SUITE_NAME")
	testCases := testcases.OpenSourceTestCasesWithRegexes(skipSuiteName, focusedSuiteName)

	boshConfig := common.BoshConfig{
		BoshURL:          mustHaveEnv("BOSH_ENVIRONMENT"),
		BoshClient:       mustHaveEnv("BOSH_CLIENT"),
		BoshClientSecret: mustHaveEnv("BOSH_CLIENT_SECRET"),
		BoshCertPath:     mustHaveEnv("BOSH_CERT_PATH"),
	}

	deploymentConfig := common.CloudFoundryConfig{
		Name:          mustHaveEnv("CF_DEPLOYMENT_NAME"),
		ApiUrl:        mustHaveEnv("CF_API_URL"),
		AdminUsername: mustHaveEnv("CF_ADMIN_USERNAME"),
		AdminPassword: mustHaveEnv("CF_ADMIN_PASSWORD"),
	}

	if containsTestCase(testCases, "cf-nfsbroker") {
		deploymentConfig.NFSServiceName = mustHaveEnv("NFS_SERVICE_NAME")
		deploymentConfig.NFSPlanName = mustHaveEnv("NFS_PLAN_NAME")
		deploymentConfig.NFSBrokerUser = os.Getenv("NFS_BROKER_USER")
		deploymentConfig.NFSBrokerPassword = os.Getenv("NFS_BROKER_PASSWORD")
		deploymentConfig.NFSBrokerUrl = os.Getenv("NFS_BROKER_URL")
	}

	configGetter := common.OSConfigGetter{
		DeploymentConfig: deploymentConfig,
		BoshConfig:       boshConfig,
	}

	runner.RunDisasterRecoveryAcceptanceTests(configGetter, testCases)
})

func containsTestCase(testCases []runner.TestCase, name string) bool {
	for _, tc := range testCases {
		if tc.Name() == name {
			return true
		}
	}

	return false
}
func mustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	if val == "" {
		panic(fmt.Sprintf("Env var %s not set\n", keyname))
	}
	return val
}
