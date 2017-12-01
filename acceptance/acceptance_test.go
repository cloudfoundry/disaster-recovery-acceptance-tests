package acceptance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("backing up Cloud Foundry", func() {
	var config runner.Config
	var testCases []runner.TestCase

	focusedSuiteName := os.Getenv("FOCUSED_SUITE_NAME")
	skipSuiteName := os.Getenv("SKIP_SUITE_NAME")
	testCases = runner.FilterTestCasesWithRegexes(testcases.OpenSourceTestCases(), skipSuiteName, focusedSuiteName)

	if os.Getenv("CONFIG") != "" {
		config = setConfigFromFile(os.Getenv("CONFIG"))
	} else {
		config = setConfigFromEnv(containsTestCase(testCases, "cf-nfsbroker"))
	}

	runner.RunDisasterRecoveryAcceptanceTests(config, testCases)
})

func setConfigFromFile(path string) runner.Config {
	configFromFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Could not load config from file: %s\n", path))
	}

	var cfConfig runner.CloudFoundryConfig
	err = json.Unmarshal(configFromFile, &cfConfig)
	if err != nil {
		panic("Could not unmarshal CloudFoundryConfig")
	}

	var boshConfig runner.BoshConfig
	err = json.Unmarshal(configFromFile, &boshConfig)
	if err != nil {
		panic("Could not unmarshal BoshConfig")
	}

	return runner.Config{
		DeploymentToBackup:  cfConfig,
		DeploymentToRestore: cfConfig,
		BoshConfig:          boshConfig,
	}
}

func setConfigFromEnv(shouldIncludeNfsBroker bool) runner.Config {
	boshConfig := runner.BoshConfig{
		BoshURL:          mustHaveEnv("BOSH_ENVIRONMENT"),
		BoshClient:       mustHaveEnv("BOSH_CLIENT"),
		BoshClientSecret: mustHaveEnv("BOSH_CLIENT_SECRET"),
		BoshCaCert:       mustHaveEnv("BOSH_CERT_PATH"),
	}
	deploymentConfig := runner.CloudFoundryConfig{
		Name:          mustHaveEnv("CF_DEPLOYMENT_NAME"),
		ApiUrl:        mustHaveEnv("CF_API_URL"),
		AdminUsername: mustHaveEnv("CF_ADMIN_USERNAME"),
		AdminPassword: mustHaveEnv("CF_ADMIN_PASSWORD"),
	}

	if shouldIncludeNfsBroker {
		deploymentConfig.NFSServiceName = mustHaveEnv("NFS_SERVICE_NAME")
		deploymentConfig.NFSPlanName = mustHaveEnv("NFS_PLAN_NAME")
		deploymentConfig.NFSBrokerUser = os.Getenv("NFS_BROKER_USER")
		deploymentConfig.NFSBrokerPassword = os.Getenv("NFS_BROKER_PASSWORD")
		deploymentConfig.NFSBrokerUrl = os.Getenv("NFS_BROKER_URL")
	}

	return runner.Config{
		DeploymentToBackup:  deploymentConfig,
		DeploymentToRestore: deploymentConfig,
		BoshConfig:          boshConfig,
	}
}

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
