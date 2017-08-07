package common

type CloudFoundryConfig struct {
	ApiUrl        string
	AdminUsername string
	AdminPassword string
}

type BoshConfig struct {
	BoshURL          string
	BoshClient       string
	BoshClientSecret string
	BoshCertPath     string
}

type Config struct {
	DeploymentToBackup  CloudFoundryConfig
	DeploymentToRestore CloudFoundryConfig
	BoshConfig          BoshConfig
}
