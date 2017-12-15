package runner

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"strings"

	"time"

	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GetAppUrl(appName string) string {
	appStats := string(RunCommandAndRetry("cf app "+appName, 5).Out.Contents())
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
	response, err := client.Get("https://" + url)
	Expect(err).NotTo(HaveOccurred())
	return response
}

func GetWithRetries(url string, retries int) *http.Response {
	var response *http.Response
	for i := 0; i < retries; i++ {
		response = Get(url)
		if response.StatusCode == http.StatusOK {
			return response
		}
		time.Sleep(10 * time.Second)
	}

	return response
}

func StatusCode(rawUrl string) func() (int, error) {
	parsedUrl, err := url.Parse(rawUrl)
	Expect(err).NotTo(HaveOccurred(), "error parsing api url")
	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
	}

	return func() (int, error) {
		client := &http.Client{
			Timeout: time.Minute,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}}
		fmt.Fprintf(GinkgoWriter, "Trying to connect to api url: %s\n", parsedUrl.String())
		resp, err := client.Get(parsedUrl.String())
		Expect(err).NotTo(HaveOccurred(), "error connecting to api url")
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
