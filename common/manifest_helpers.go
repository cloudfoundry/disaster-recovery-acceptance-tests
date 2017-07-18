package common

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	"gopkg.in/yaml.v1"
)

func FindCredentialsFor(deploymentName string) (string, string, string) {
	manifest := DownloadManifest(deploymentName)
	deployment := Deployment{}
	yaml.Unmarshal([]byte(manifest), &deployment)
	domain := deployment.FindAppDomain()
	username, password := deployment.FindCredentials()
	return fmt.Sprintf("https://api.%s", domain), username, password
}

func (deployment Deployment) FindAppDomain() string {
	for _, instanceGroup := range deployment.InstanceGroups {
		if instanceGroup.Name == "api" {
			for _, job := range instanceGroup.Jobs {
				if job.Name == "cloud_controller_ng" {
					return job.Properties.AppDomains[0]
				}
			}
		}
	}

	Fail("No app domains found")
	return ""
}

func (deployment Deployment) FindCredentials() (string, string) {
	for _, instanceGroup := range deployment.InstanceGroups {
		if instanceGroup.Name == "uaa" {
			for _, job := range instanceGroup.Jobs {
				if job.Name == "uaa" {
					user := job.Properties.UAA.SCIM.Users[0]
					return user.Name, user.Password
				}
			}
		}
	}

	Fail("No admin credentials found")
	return "", ""
}

type Deployment struct {
	InstanceGroups []struct {
		Name string
		Jobs []struct {
			Name       string
			Properties struct {
				AppDomains []string `yaml:"app_domains"`
				UAA        struct {
					SCIM struct {
						Users []struct {
							Name     string
							Password string
						}
					} `yaml:"scim"`
				} `yaml:"uaa"`
			}
		}
	} `yaml:"instance_groups"`
}
