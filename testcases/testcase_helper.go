package testcases

import (
	"fmt"
	"regexp"

	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
)

func OpenSourceTestCases() []runner.TestCase {
	return []runner.TestCase{
		NewAppUptimeTestCase(),
		NewCfAppTestCase(),
		NewCfUaaTestCase(),
		NewCfNetworkingTestCase(),
		NewNFSTestCases(),
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

// Return all test cases whose names DO NOT match the skip regex or DO match the focus regex
// (Inspired by ginkgo's --skip and --focus flags)
func OpenSourceTestCasesWithRegexes(skip, focus string) []runner.TestCase {
	allCases := OpenSourceTestCases()
	if skip == "" && focus == "" {
		return allCases
	}

	testCases := []runner.TestCase{}
	for _, tc := range allCases {
		if !shouldSkipCase(skip, focus, tc) {
			fmt.Printf("Running %s test case\n", tc.Name())
			testCases = append(testCases, tc)
		}
	}

	if (len(testCases)) > 0 {
		return testCases
	}

	panic(fmt.Sprintf("Unable to find test case matching regex %s for focus or %s for skipping\n", focus, skip))
}

func shouldSkipCase(skip, focus string, tc runner.TestCase) bool {
	matchesFocus := true
	matchesSkip := false

	caseName := tc.Name()

	if focus != "" {
		focusFilter := regexp.MustCompile(focus)
		matchesFocus = focusFilter.MatchString(caseName)
	}

	if skip != "" {
		skipFilter := regexp.MustCompile(skip)
		matchesSkip = skipFilter.MatchString(caseName)
	}

	return !matchesFocus || matchesSkip
}
