package acceptance

import (
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("backing up Cloud Foundry", func() {
	configGetter := NewOSConfigGetter()

	runner.RunDisasterRecoveryAcceptanceTests(configGetter, testcases.OpenSourceTestCases())
})

