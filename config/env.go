package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
)

func FromEnv() (runner.Config, runner.TestCaseFilter) {
	boshConfig := runner.BoshConfig{
		BoshURL:          mustHaveEnv("BOSH_ENVIRONMENT"),
		BoshClient:       mustHaveEnv("BOSH_CLIENT"),
		BoshClientSecret: mustHaveEnv("BOSH_CLIENT_SECRET"),
		BoshCaCert:       mustHaveEnv("BOSH_CA_CERT"),
	}
	deploymentConfig := runner.CloudFoundryConfig{
		Name:          mustHaveEnv("CF_DEPLOYMENT_NAME"),
		ApiUrl:        mustHaveEnv("CF_API_URL"),
		AdminUsername: mustHaveEnv("CF_ADMIN_USERNAME"),
		AdminPassword: mustHaveEnv("CF_ADMIN_PASSWORD"),
	}

	deploymentConfig.NFSServiceName = os.Getenv("NFS_SERVICE_NAME")
	deploymentConfig.NFSPlanName = os.Getenv("NFS_PLAN_NAME")
	deploymentConfig.NFSBrokerUser = os.Getenv("NFS_BROKER_USER")
	deploymentConfig.NFSBrokerPassword = os.Getenv("NFS_BROKER_PASSWORD")
	deploymentConfig.NFSBrokerUrl = os.Getenv("NFS_BROKER_URL")

	deploymentConfig.SMBServiceName = os.Getenv("SMB_SERVICE_NAME")
	deploymentConfig.SMBPlanName = os.Getenv("SMB_PLAN_NAME")
	deploymentConfig.SMBBrokerUser = os.Getenv("SMB_BROKER_USER")
	deploymentConfig.SMBBrokerPassword = os.Getenv("SMB_BROKER_PASSWORD")
	deploymentConfig.SMBBrokerUrl = os.Getenv("SMB_BROKER_URL")

	timeout := TimeoutFromEnv()

	deleteAndRedeployCF := os.Getenv("DELETE_AND_REDEPLOY_CF") == "true"

	conf := runner.Config{
		CloudFoundryConfig:  deploymentConfig,
		BoshConfig:          boshConfig,
		Timeout:             timeout,
		DeleteAndRedeployCF: deleteAndRedeployCF,
	}

	filter := runner.NewRegexTestCaseFilter(os.Getenv("FOCUSED_SUITE_NAME"), os.Getenv("SKIP_SUITE_NAME"))

	return conf, filter
}

func TimeoutFromEnv() time.Duration {
	var timeout time.Duration
	timeoutStr := os.Getenv("DEFAULT_TIMEOUT_MINS")
	if timeoutStr != "" {
		timeoutInt, err := strconv.Atoi(timeoutStr)
		if err != nil {
			panic(fmt.Sprint("DEFAULT_TIMEOUT_MINS, if set, must be an integer\n"))
		}
		timeout = time.Duration(timeoutInt) * time.Minute
	} else {
		timeout = defaultTimeout
	}

	return timeout
}

func mustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	if val == "" {
		panic(fmt.Sprintf("Env var %s not set\n", keyname))
	}
	return val
}
