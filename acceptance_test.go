package disaster_recovery_acceptance_tests

import (
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/testcases"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("backing up Cloud Foundry", func() {
	runner.RunDisasterRecoveryAcceptanceTests([]runner.TestCase{
		testcases.NewAppUptimeTestCase(),
		testcases.NewCfAppTestCase(),
	})
})
