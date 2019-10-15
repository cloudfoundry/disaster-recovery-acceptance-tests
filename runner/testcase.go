package runner

type TestCase interface {
	Name() string
	CheckDeployment(Config)
	BeforeBackup(Config)
	AfterBackup(Config)
	EnsureAfterSelectiveRestore(Config)
	AfterRestore(Config)
	Cleanup(Config)
}
