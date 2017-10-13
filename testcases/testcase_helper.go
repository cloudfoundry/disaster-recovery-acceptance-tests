package testcases

import (
	"fmt"

	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
)

func OpenSourceTestCases() []runner.TestCase {
	return []runner.TestCase{
		NewAppUptimeTestCase(),
		NewCfAppTestCase(),
		NewCfUaaTestCase(),
		NewCfNetworkingTestCase(),
	}
}

func OpenSourceTestCasesWithFocus(suiteName string) []runner.TestCase {
	for _, tc := range OpenSourceTestCases() {
		if tc.Name() == suiteName {
			fmt.Printf("Test focused: running %s test case\n", suiteName)
			return []runner.TestCase{tc}
		}
	}

	panic(fmt.Sprintf("Unable to find test case with name %s\n", suiteName))
}
