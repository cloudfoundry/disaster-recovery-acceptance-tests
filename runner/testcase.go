package runner

type TestCase interface {
	BeforeBackup()
	AfterBackup()
	AfterRestore()
	Cleanup()
}
