package backup_and_restore

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"

	"fmt"

	acceptance "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/acceptance/backup_and_restore/test_cases"

	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

func TestPcfBackupAndRestoreAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PcfBackupAndRestoreAcceptanceTests Suite")
}

var testCases []acceptance.TestCase
var uniqueTestID string
var jumpBoxSession *JumpBoxSession

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)
	uniqueTestID = RandomStringNumber()
	jumpBoxSession = NewJumpBoxSession(uniqueTestID)

	// ### test cases to be run
	testCases = []acceptance.TestCase{
		acceptance.NewAppUptimeTestCase(),
		acceptance.NewCfAppTestCase(),
	}
})

var _ = AfterSuite(func() {
	jumpBoxSession.Cleanup()
})

func printEnvsAreDifferentWarning() {
	fmt.Println("     --------------------------------------------------------")
	fmt.Println("     NOTE: this suite is currently configured to back up from")
	fmt.Println("     one environment and restore to a difference one. Make   ")
	fmt.Println("     sure this is the intended configuration.                ")
	fmt.Println("     --------------------------------------------------------")
}
