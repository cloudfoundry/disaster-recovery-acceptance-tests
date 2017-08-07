package runner

import "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"

type TestCase interface {
	BeforeBackup(common.Config)
	AfterBackup(common.Config)
	AfterRestore(common.Config)
	Cleanup(common.Config)
}
