package runner

type CloudFoundryConfig struct {
	Name                              string
	ApiUrl                            string
	AdminUsername                     string
	AdminPassword                     string
	NotificationsTemplateClientID     string
	NotificationsTemplateClientSecret string
	NFSServiceName                    string
	NFSPlanName                       string
	NFSBrokerUser                     string
	NFSBrokerPassword                 string
	NFSBrokerUrl                      string
}

type BoshConfig struct {
	BoshURL          string
	BoshClient       string
	BoshClientSecret string
	BoshCertPath     string
}

type Config struct {
	Deployment CloudFoundryConfig
	BoshConfig BoshConfig
}

type ConfigGetter interface {
	FindConfig([]TestCase) Config
}