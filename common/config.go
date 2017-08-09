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
	manifestData, err := unstructured.ParseYAML(DownloadManifest(deploymentName, configGetter.BoshConfig))
	if err != nil {
		panic("Error downloading manifest")
	}

	instanceGroups, err := manifestData.GetByPointer("/instance_groups")
	if err != nil {
		panic("Error parsing manifest")
	}

	apiGroup, found := instanceGroups.FindElem(ByName("api"))
	if !found {
		panic("Error parsing manifest")
	}
	apiGroupJobs, err := apiGroup.GetByPointer("/jobs")
	if err != nil {
		panic("Error parsing manifest")
	}
	cloudControllerJob, found := apiGroupJobs.FindElem(ByName("cloud_controller_ng"))
	if !found {
		panic("Error parsing manifest")
	}
	systemDomain, err := cloudControllerJob.GetByPointer("/properties/system_domain")
	if err != nil {
		panic("Error parsing manifest")
	}

	uaaGroup, found := instanceGroups.FindElem(ByName("uaa"))
	if !found {
		panic("Error parsing manifest")
	}
	uaaGroupJobs, err := uaaGroup.GetByPointer("/jobs")
	if err != nil {
		panic("Error parsing manifest")
	}
	uaaJob, found := uaaGroupJobs.FindElem(ByName("uaa"))
	if !found {
		panic("Error parsing manifest")
	}
	user, err := uaaJob.GetByPointer("/properties/uaa/scim/users/0")
	if err != nil {
		panic("Error parsing manifest")
	}
	username, err := user.GetByPointer("/name")
	if err != nil {
		panic("Error parsing manifest")
	}
	password, err := user.GetByPointer("/password")
	if err != nil {
		panic("Error parsing manifest")
	}

	return CloudFoundryConfig{
		Name:          deploymentName,
		ApiUrl:        fmt.Sprintf("https://api.%s", systemDomain.UnsafeStringValue()),
		AdminUsername: username.UnsafeStringValue(),
		AdminPassword: password.UnsafeStringValue(),
	}
}

func ByName(name string) unstructured.ElementMatcher {
	return func(element unstructured.Data) bool {
		return element.HasKey("name") &&
			element.F("name").UnsafeStringValue() == name
	}
}
