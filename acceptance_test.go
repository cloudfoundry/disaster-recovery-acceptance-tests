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
	boshConfig := common.BoshConfig{
		BoshURL:          mustHaveEnv("BOSH_ENVIRONMENT"),
		BoshClient:       mustHaveEnv("BOSH_CLIENT"),
		BoshClientSecret: mustHaveEnv("BOSH_CLIENT_SECRET"),
		BoshCertPath:     mustHaveEnv("BOSH_CERT_PATH"),
	}

	deploymentConfig := common.CloudFoundryConfig{
		Name:              mustHaveEnv("CF_DEPLOYMENT_NAME"),
		ApiUrl:            mustHaveEnv("CF_API_URL"),
		AdminUsername:     mustHaveEnv("CF_ADMIN_USERNAME"),
		AdminPassword:     mustHaveEnv("CF_ADMIN_PASSWORD"),
		NFSServiceName:    mustHaveEnv("NFS_SERVICE_NAME"),
		NFSPlanName:       mustHaveEnv("NFS_PLAN_NAME"),
		NFSBrokerUser:     os.Getenv("NFS_BROKER_USER"),
		NFSBrokerPassword: os.Getenv("NFS_BROKER_PASSWORD"),
		NFSBrokerUrl:      os.Getenv("NFS_BROKER_URL"),
	}

	configGetter := common.OSConfigGetter{
		DeploymentConfig: deploymentConfig,
		BoshConfig:       boshConfig,
	}

	var testCases []runner.TestCase

	focusedSuiteName := os.Getenv("FOCUSED_SUITE_NAME")
	if focusedSuiteName != "" {
		testCases = testcases.OpenSourceTestCasesWithFocus(focusedSuiteName)
	} else {
		testCases = testcases.OpenSourceTestCases()
	}

	runner.RunDisasterRecoveryAcceptanceTests(configGetter, testCases)
})

func mustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	if val == "" {
		panic(fmt.Sprintf("Env var %s not set\n", keyname))
	}
	return val
}
