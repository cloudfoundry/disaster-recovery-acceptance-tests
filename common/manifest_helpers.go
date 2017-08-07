package common

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	"gopkg.in/yaml.v1"
)

func FindCredentialsFor(deploymentName string, boshConfig BoshConfig) (string, string, string) {
	manifest := DownloadManifest(deploymentName, boshConfig)
	deployment := Deployment{}
	yaml.Unmarshal([]byte(manifest), &deployment)
	systemDomain := deployment.FindSystemDomain()
	username, password := deployment.FindCredentials()
	return fmt.Sprintf("https://api.%s", systemDomain), username, password
}

func (deployment Deployment) FindSystemDomain() string {
	for _, instanceGroup := range deployment.InstanceGroups {
		if instanceGroup.Name == "api" || //OS
			instanceGroup.Name == "cloud_controller" { //ERT
			for _, job := range instanceGroup.Jobs {
				if job.Name == "cloud_controller_ng" {
					if job.Properties.SystemDomain != "" {
						return job.Properties.SystemDomain
					} else if instanceGroup.Properties.SystemDomain != "" {
						return instanceGroup.Properties.SystemDomain

					} else {
						Fail("CC job found, but no system domain found")
					}
				}
			}
		}
	}

	Fail("No system domain found")
	return ""
}

func (deployment Deployment) FindCredentials() (string, string) {
	for _, instanceGroup := range deployment.InstanceGroups {
		if instanceGroup.Name == "uaa" {
			for _, job := range instanceGroup.Jobs {
				if job.Name == "uaa" {
					if len(job.Properties.UAA.SCIM.Users) > 0 {
						user := job.Properties.UAA.SCIM.Users[0] //OS
						return user.Name, user.Password
					} else if len(instanceGroup.Properties.UAA.SCIM.Users) > 0 {
						user := instanceGroup.Properties.UAA.SCIM.Users[0] //ERT
						return user.Name, user.Password
					} else {
						Fail("UAA job found, but no admin user found")
					}
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
				SystemDomain string `yaml:"system_domain"`
				UAA          struct {
					SCIM struct {
						Users []struct {
							Name     string
							Password string
						}
					} `yaml:"scim"`
				} `yaml:"uaa"`
			}
		}
		Properties struct {
			SystemDomain string `yaml:"system_domain"`
			UAA          struct {
				SCIM struct {
					Users []struct {
						Name     string
						Password string
					}
				} `yaml:"scim"`
			} `yaml:"uaa"`
		}
	} `yaml:"instance_groups"`
}
