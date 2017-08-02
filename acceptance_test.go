package disaster_recovery_acceptance_tests

import (
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("backing up Cloud Foundry", func() {
	config := ReadConfigFromBOSHManifest()
	runner.RunDisasterRecoveryAcceptanceTests(config, []runner.TestCase{
		testcases.NewAppUptimeTestCase(),
		testcases.NewCfAppTestCase(),
	})
})

func ReadConfigFromBOSHManifest() runner.Config {
	urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())
	urlForDeploymentToRestore, usernameForDeploymentToRestore, passwordForDeploymentToRestore := FindCredentialsFor(DeploymentToRestore())

	return runner.Config{
		DeploymentToBackup: runner.CloudFoundryConfig{
			ApiUrl:        urlForDeploymentToBackup,
			AdminUsername: usernameForDeploymentToBackup,
			AdminPassword: passwordForDeploymentToBackup,
		},
		DeploymentToRestore: runner.CloudFoundryConfig{
			ApiUrl:        urlForDeploymentToRestore,
			AdminUsername: usernameForDeploymentToRestore,
			AdminPassword: passwordForDeploymentToRestore,
		},
	}

}
