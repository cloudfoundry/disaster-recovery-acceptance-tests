package common

type CloudFoundryConfig struct {
	Name                              string
	ApiUrl                            string
	AdminUsername                     string
	AdminPassword                     string
	NotificationsTemplateClientID     string
	NotificationsTemplateClientSecret string
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

type ConfigGetter interface {
	FindConfig() Config
}

type OSConfigGetter struct {
	DeploymentConfig CloudFoundryConfig
	BoshConfig       BoshConfig
}

func (configGetter OSConfigGetter) FindConfig() Config {
	return Config{
		DeploymentToBackup:  configGetter.DeploymentConfig,
		DeploymentToRestore: configGetter.DeploymentConfig,
		BoshConfig:          configGetter.BoshConfig,
	}
}
