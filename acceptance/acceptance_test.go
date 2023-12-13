package acceptance

import (
	"os"

	"github.com/cloudfoundry/disaster-recovery-acceptance-tests/config"
	"github.com/cloudfoundry/disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry/disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("backing up Cloud Foundry", func() {
	var conf runner.Config
	var filter runner.TestCaseFilter

	if os.Getenv("CONFIG") != "" {
		conf, filter = config.FromFile(os.Getenv("CONFIG"))
	} else {
		conf, filter = config.FromEnv()
	}

  config.Timeout = "60m"
	runner.RunDisasterRecoveryAcceptanceTests(conf, filter.Filter(testcases.OpenSourceTestCases()))
})
