package runner

type CloudFoundryConfig struct {
	ApiUrl        string
	AdminUsername string
	AdminPassword string
}

type Config struct {
	DeploymentToBackup  CloudFoundryConfig
	DeploymentToRestore CloudFoundryConfig
}
