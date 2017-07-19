package acceptance

type TestCase interface {
	PopulateState()
	CheckState()
	Cleanup()
}
