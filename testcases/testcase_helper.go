package testcases

import (
	"github.com/cloudfoundry/disaster-recovery-acceptance-tests/runner"
)

func OpenSourceTestCases() []runner.TestCase {
	return []runner.TestCase{
		NewRouterGroupTestCase(),
		NewAppUptimeTestCase(),
		NewCfAppTestCase(),
		NewCfUaaTestCase(),
		NewCfNetworkingTestCase(),
		NewNFSTestCases(),
		NewSMBTestCases(),
		NewCfCredhubSSITestCase(),
	}
}
