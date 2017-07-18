package common

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GetAppUrl(appName string) string {
	appStats := string(RunCommandSuccessfully("cf app " + appName).Out.Contents())
	var appUrl string
	for _, line := range strings.Split(appStats, "\n") {
		if strings.HasPrefix(line, "routes:") {
			s := strings.Split(line, " ")
			appUrl = s[len(s)-1]
		}
	}

	Expect(appUrl).NotTo(BeEmpty())
	return appUrl
}
func Get(url string) *http.Response {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	fmt.Fprintf(GinkgoWriter, "Polling url %s\n", url)
	response, err := client.Get("https://" + url)
	fmt.Fprintf(GinkgoWriter, "Done polling url %s\n", url)
	Expect(err).NotTo(HaveOccurred())
	return response
}

func GetCurrentApplicationStateFor(guid string) (string, error) {
	statusJson := RunCommandSuccessfully(fmt.Sprintf("cf curl /v2/apps/%s/stats", string(guid)))
	response := AppStatusResponse{}
	unmarshalError := json.Unmarshal(statusJson.Out.Contents(), &response)

	if unmarshalError != nil {
		return "", unmarshalError
	}

	return response["0"].State, nil
}

type InstanceStatusResponse struct {
	State string
}

type AppStatusResponse map[string]InstanceStatusResponse
