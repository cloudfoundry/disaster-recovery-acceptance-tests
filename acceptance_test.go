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
		BoshURL:          mustHaveEnv("BOSH_URL"),
		BoshClient:       mustHaveEnv("BOSH_CLIENT"),
		BoshClientSecret: mustHaveEnv("BOSH_CLIENT_SECRET"),
		BoshCertPath:     mustHaveEnv("BOSH_CERT_PATH"),
	}

	configGetter := common.OSConfigGetter{
		DeploymentNameForBackup:  mustHaveEnv("DEPLOYMENT_TO_BACKUP"),
		DeploymentNameForRestore: mustHaveEnv("DEPLOYMENT_TO_RESTORE"),
		BoshConfig:               boshConfig,
	}

	runner.RunDisasterRecoveryAcceptanceTests(configGetter, testcases.OpenSourceTestCases())
})

func mustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	if val == "" {
		panic(fmt.Sprintln("Env var %s not set", keyname))
	}
	return val
}
