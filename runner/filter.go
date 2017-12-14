package runner

import (
	"fmt"
	"regexp"
	"strings"
)

type TestCaseFilter interface {
	Filter([]TestCase) []TestCase
}

type RegexTestCaseFilter struct {
	focusedSuiteName string
	skipSuiteName    string
}

func NewRegexTestCaseFilter(focusedSuiteName, skipSuiteName string) RegexTestCaseFilter {
	return RegexTestCaseFilter{focusedSuiteName: focusedSuiteName, skipSuiteName: skipSuiteName}
}

func (f RegexTestCaseFilter) Filter(testCases []TestCase) []TestCase {
	if f.skipSuiteName == "" && f.focusedSuiteName == "" {
		return testCases
	}

	var filteredTestCases []TestCase
	for _, testCase := range testCases {
		if shouldRun(f.skipSuiteName, f.focusedSuiteName, testCase) {
			filteredTestCases = append(filteredTestCases, testCase)
		}
	}

	if (len(filteredTestCases)) > 0 {
		return filteredTestCases
	}

	panic(fmt.Sprintf("Unable to find test case matching regex %s for focus or %s for skipping\n", f.focusedSuiteName, f.skipSuiteName))
}

func shouldRun(skip, focus string, testCase TestCase) bool {
	caseName := testCase.Name()

	matchesFocus := true
	if focus != "" {
		focusFilter := regexp.MustCompile(focus)
		matchesFocus = focusFilter.MatchString(caseName)
	}

	matchesSkip := false
	if skip != "" {
		skip = strings.TrimSpace(skip)
		skipFilter := regexp.MustCompile(skip)
		matchesSkip = skipFilter.MatchString(caseName)
	}

	return matchesFocus && !matchesSkip
}

type IntegrationConfigTestCaseFilter map[string]interface{}

func (f IntegrationConfigTestCaseFilter) Filter(testCases []TestCase) []TestCase {
	var filteredTestCases []TestCase
	for _, testCase := range testCases {
		if f["include_"+testCase.Name()] == true {
			filteredTestCases = append(filteredTestCases, testCase)
		}
	}

	if (len(filteredTestCases)) > 0 {
		return filteredTestCases
	}

	panic(fmt.Sprintf("Unable to find any test case included by the config"))
}
