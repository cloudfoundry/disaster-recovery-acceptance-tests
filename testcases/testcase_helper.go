package testcases

import (
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
)

func OpenSourceTestCases() []runner.TestCase {
	return []runner.TestCase{
		NewRouterGroupTestCase(),
		NewAppUptimeTestCase(),
		NewCfAppTestCase(),
		NewCfUaaTestCase(),
		NewCfNetworkingTestCase(),
		NewNFSTestCases(),
		NewCfCredhubSSITestCase(),
	}
}

func ExperimentalTestCases() []runner.TestCase {
	return []runner.TestCase{
		NewPermTestCase(),
	}
}
