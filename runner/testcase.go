package runner

type TestCase interface {
	Name() string
	BeforeBackup(Config)
	AfterBackup(Config)
	AfterRestore(Config)
	Cleanup(Config)
}
