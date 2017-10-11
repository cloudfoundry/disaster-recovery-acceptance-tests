package testcases

import "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"

func OpenSourceTestCases() []runner.TestCase {
	return []runner.TestCase{
		NewAppUptimeTestCase(),
		NewCfAppTestCase(),
		NewCfUaaTestCase(),
	}
	
}
