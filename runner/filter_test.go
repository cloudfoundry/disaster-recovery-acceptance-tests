package runner_test

import (
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestCaseFilter", func() {
	Describe("RegexTestCaseFilter", func() {
		var filter runner.RegexTestCaseFilter
		var focusRegex, skipRegex string

		JustBeforeEach(func() {
			filter = runner.NewRegexTestCaseFilter(focusRegex, skipRegex)
		})

		Context("when neither regex is specified", func() {
			BeforeEach(func() {
				focusRegex = ""
				skipRegex = ""
			})

			It("includes all tests", func() {
				Expect(filter.Filter(testCases("one", "two", "three"))).To(ConsistOf(
					testCase("one"),
					testCase("two"),
					testCase("three"),
				))
			})
		})

		Context("when only the focus regex is specified", func() {
			BeforeEach(func() {
				focusRegex = "one"
				skipRegex = ""
			})

			It("only includes the tests matching the focus regex", func() {
				Expect(filter.Filter(testCases("one", "two", "three"))).To(ConsistOf(
					testCase("one"),
				))
			})
		})

		Context("when only the skip regex is specified", func() {
			BeforeEach(func() {
				focusRegex = ""
				skipRegex = "one"
			})

			It("excludes the tests matching the skip regex", func() {
				Expect(filter.Filter(testCases("one", "two", "three"))).To(ConsistOf(
					testCase("two"),
					testCase("three"),
				))
			})
		})

		Context("when both the focus regex and the skip regex are specified", func() {
			BeforeEach(func() {
				focusRegex = "one|three"
				skipRegex = "two|three"
			})

			It("only runs tests that match the focus regex but not the skip regex", func() {
				Expect(filter.Filter(testCases("one", "two", "three", "four"))).To(ConsistOf(
					testCase("one"),
				))
			})
		})

		Context("when no test case matches", func() {
			BeforeEach(func() {
				focusRegex = "foo"
			})

			It("panics", func() {
				Expect(func() {
					filter.Filter(testCases("bar"))
				}).To(Panic())
			})
		})
	})

	Describe("IntegrationConfigTestCaseFilter", func() {
		var filter runner.IntegrationConfigTestCaseFilter

		JustBeforeEach(func() {
			filter = runner.IntegrationConfigTestCaseFilter(map[string]interface{}{
				"include_one":  true,
				"include_two":  false,
				"include_four": "some value",
				"include_five": true,
			})
		})

		It("only runs tests that are included in the config", func() {
			Expect(filter.Filter(testCases("one", "two", "three", "four"))).To(ConsistOf(
				testCase("one"),
			))
		})

		Context("when no test case matches", func() {
			It("panics", func() {
				Expect(func() {
					filter.Filter(testCases("six"))
				}).To(Panic())
			})
		})
	})
})

type FakeTestCase struct {
	name string
}

func (tc FakeTestCase) Name() string {
	return tc.name
}

func (tc FakeTestCase) BeforeBackup(config runner.Config) {}
func (tc FakeTestCase) AfterBackup(config runner.Config)  {}
func (tc FakeTestCase) AfterRestore(config runner.Config) {}
func (tc FakeTestCase) Cleanup(config runner.Config)      {}

func testCase(name string) runner.TestCase {
	return FakeTestCase{name: name}
}

func testCases(names ...string) []runner.TestCase {
	var tcs []runner.TestCase

	for _, name := range names {
		tcs = append(tcs, testCase(name))
	}

	return tcs
}
