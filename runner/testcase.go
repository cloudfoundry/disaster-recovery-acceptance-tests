package runner

type TestCase interface {
	Name() string
	CheckDeployment(Config)
	BeforeBackup(Config)
	AfterBackup(Config)
	AfterRestore(Config)
	Cleanup(Config)
}
