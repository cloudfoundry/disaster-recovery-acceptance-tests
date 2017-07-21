package backup_and_restore

type TestCase interface {
	BeforeBackup()
	AfterBackup()
	AfterRestore()
	Cleanup()
}
