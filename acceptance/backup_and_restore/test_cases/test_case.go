package backup_and_restore

type TestCase interface {
	PopulateState()
	CheckState()
	Cleanup()
}
