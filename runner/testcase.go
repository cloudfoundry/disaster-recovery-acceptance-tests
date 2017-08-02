package runner

type TestCase interface {
	BeforeBackup(Config)
	AfterBackup(Config)
	AfterRestore(Config)
	Cleanup(Config)
}
