package common

import (
	"fmt"

	"github.com/totherme/unstructured"
)

type CloudFoundryConfig struct {
	Name          string
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

type ConfigGetter interface {
	FindConfig() Config
}

type OSConfigGetter struct {
	DeploymentNameForBackup  string
	DeploymentNameForRestore string
	BoshConfig               BoshConfig
}

func (configGetter OSConfigGetter) FindConfig() Config {
	return Config{
		DeploymentToBackup:  configGetter.findCloudFoundryConfigFor(configGetter.DeploymentNameForBackup),
		DeploymentToRestore: configGetter.findCloudFoundryConfigFor(configGetter.DeploymentNameForRestore),
		BoshConfig:          configGetter.BoshConfig,
	}
}

func (configGetter OSConfigGetter) findCloudFoundryConfigFor(deploymentName string) CloudFoundryConfig {
	// TODO: deal with errors, make helpful parsing error logs

	manifestData, _ := unstructured.ParseYAML(DownloadManifest(deploymentName, configGetter.BoshConfig))

	instanceGroups, _ := manifestData.GetByPointer("/instance_groups")

	apiGroup, _ := instanceGroups.FindElem(ByName("api"))
	apiGroupJobs, _ := apiGroup.GetByPointer("/jobs")
	cloudControllerJob, _ := apiGroupJobs.FindElem(ByName("cloud_controller_ng"))
	systemDomain, _ := cloudControllerJob.GetByPointer("/properties/system_domain")

	uaaGroup, _ := instanceGroups.FindElem(ByName("uaa"))
	uaaGroupJobs, _ := uaaGroup.GetByPointer("/jobs")
	uaaJob, _ := uaaGroupJobs.FindElem(ByName("uaa"))
	user, _ := uaaJob.GetByPointer("/properties/uaa/scim/users/0")
	username, _ := user.GetByPointer("/name")
	password, _ := user.GetByPointer("/password")

	return CloudFoundryConfig{
		Name:          deploymentName,
		ApiUrl:        fmt.Sprintf("https://api.%s", systemDomain.UnsafeStringValue()),
		AdminUsername: username.UnsafeStringValue(),
		AdminPassword: password.UnsafeStringValue(),
	}
}

//TODO: is there a better place for this?
func ByName(name string) unstructured.ElementMatcher {
	return func(element unstructured.Data) bool {
		return element.HasKey("name") &&
			element.F("name").UnsafeStringValue() == name
	}
}
