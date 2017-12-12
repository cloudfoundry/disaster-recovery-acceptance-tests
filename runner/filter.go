package runner

import (
	"regexp"
	"strings"
	"fmt"
)

// Return all test cases whose names DO NOT match the skip regex or DO match the focus regex
// (Inspired by ginkgo's --skip and --focus flags)
func FilterTestCasesWithRegexes(allCases []TestCase, skip, focus string) []TestCase {
	if skip == "" && focus == "" {
		return allCases
	}

	testCases := []TestCase{}
	for _, tc := range allCases {
		if !shouldSkipCase(skip, focus, tc) {
			testCases = append(testCases, tc)
		}
	}

	if (len(testCases)) > 0 {
		return testCases
	}

	panic(fmt.Sprintf("Unable to find test case matching regex %s for focus or %s for skipping\n", focus, skip))
}

func shouldSkipCase(skip, focus string, tc TestCase) bool {
	matchesFocus := true
	matchesSkip := false

	caseName := tc.Name()

	if focus != "" {
		focusFilter := regexp.MustCompile(focus)
		matchesFocus = focusFilter.MatchString(caseName)
	}

	if skip != "" {
		skip = strings.TrimSpace(skip)
		skipFilter := regexp.MustCompile(skip)
		matchesSkip = skipFilter.MatchString(caseName)
	}

	return !matchesFocus || matchesSkip
}

