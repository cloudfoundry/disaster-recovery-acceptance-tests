package common

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"crypto/rand"

	"encoding/base64"

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

func StatusCode(url string) func() (int, error) {
	return func() (int, error) {
		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
		resp, err := client.Get(url)
		if err != nil {
			return 0, err
		}
		return resp.StatusCode, err
	}
}

func RandomStringNumber() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random string number")
	}

	return base64.RawURLEncoding.EncodeToString(b)
}

type InstanceStatusResponse struct {
	State string
}

type AppStatusResponse map[string]InstanceStatusResponse
