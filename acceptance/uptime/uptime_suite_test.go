package uptime

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strconv"
	"testing"
	"time"

	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

func TestAppUptime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AppUptime Suite")
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
