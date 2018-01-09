package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
)

const defaultTimeout = 15 * time.Minute

func FromFile(path string) (runner.Config, runner.TestCaseFilter) {
	configFromFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprint(fmt.Sprintf("Could not load config from file: %s\n", path)))
	}

	conf := runner.Config{}
	err = json.Unmarshal(configFromFile, &conf)
	if err != nil {
		panic(fmt.Sprint("Could not unmarshal Config"))
	}

	timeoutConfig := timeoutConfig{}
	err = json.Unmarshal(configFromFile, &timeoutConfig)
	if err != nil {
		panic(fmt.Sprint("Could not unmarshal timeout"))
	}

	if timeoutConfig.TimeoutInMinutes == 0 {
		conf.Timeout = defaultTimeout
	} else {
		conf.Timeout = time.Minute * time.Duration(timeoutConfig.TimeoutInMinutes)
	}

	filter := runner.IntegrationConfigTestCaseFilter{}
	err = json.Unmarshal(configFromFile, &filter)
	if err != nil {
		panic(fmt.Sprint("Could not unmarshal Filter"))
	}

	return conf, filter
}

type timeoutConfig struct {
	TimeoutInMinutes int64 `json:"timeout_in_minutes"`
}
