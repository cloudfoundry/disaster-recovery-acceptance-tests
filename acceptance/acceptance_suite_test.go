package acceptance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strconv"
	"testing"
	"time"

	"github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/acceptance/testcases"
	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

func TestPcfBackupAndRestoreAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PcfBackupAndRestoreAcceptanceTests Suite")
}

// ### test cases to be run
var testCases = []acceptance.TestCase{
	acceptance.NewCfAppTestCase(uniqueTestID),
}

var uniqueTestID string
var jumpBoxSession *JumpBoxSession

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)
	uniqueTestID = timestamp()
	jumpBoxSession = NewJumpBoxSession(uniqueTestID)
})

var _ = AfterSuite(func() {
	jumpBoxSession.Cleanup()
})

func timestamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 16)
}
