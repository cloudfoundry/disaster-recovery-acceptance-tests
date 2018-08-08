package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"encoding/json"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
)

const credhubBaseURL = "https://credhub.service.cf.internal:8844"
const credentialName = "SECRET_PASSWORD"

func main() {
	s, err := newServer()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/permission", s.Permission)
	http.HandleFunc("/create", s.Create)
	http.HandleFunc("/list", s.List)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

type Server struct {
	client  *http.Client
	counter int
}

type CredhubClient struct {
	CredhubClient string `json:"credhub_client"`
	CredhubSecret string `json:"credhub_secret"`
}

type VcapApplication struct {
	ApplicationId string `json:"application_id"`
}

func (s *Server) Permission(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered error reading body: [%s]", err.Error())
		return
	}

	var creds CredhubClient
	err = json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered error getting client creds from permission request: [%s]", err.Error())
		return
	}

	ch, err := credhub.New(
		credhubBaseURL,
		credhub.SkipTLSValidation(true),
		credhub.Auth(auth.UaaClientCredentials(creds.CredhubClient, creds.CredhubSecret)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating credhub client: [%s]", err.Error())
		return
	}

	var vcapApp VcapApplication
	err = json.Unmarshal([]byte(os.Getenv("VCAP_APPLICATION")), &vcapApp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error unmarshalling vcap_application: [%s]", err.Error())
		return
	}

	permissionBody := map[string]interface{}{
		"actor":      "mtls-app:" + vcapApp.ApplicationId,
		"operations": []string{"read", "write"},
		"path":       credentialName + "/*",
	}

	_, err = ch.Request("POST", "/api/v2/permissions", nil, permissionBody, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating permission in credhub: [%s]", err.Error())
		return
	}
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	s.counter++
	resp, err := s.client.Post(fmt.Sprintf("%s/api/v1/data", credhubBaseURL), "application/json",
		strings.NewReader(fmt.Sprintf(`{"name": "/%s/%d", "type": "password"}`, credentialName, s.counter)),
	)
	if ok := handleBadResponses(w, resp, err); !ok {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	resp, err := s.client.Get(fmt.Sprintf("%s/api/v1/data?path=%s", credhubBaseURL, credentialName))
	if ok := handleBadResponses(w, resp, err); !ok {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered error reading response body from credhub: [%s]", err)
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, string(body))
}

func handleBadResponses(w http.ResponseWriter, resp *http.Response, err error) bool {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered error from credhub: [%s]", err)
		return false
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintf(w, "Encountered bad response code from credhub: %d", resp.StatusCode)
		return false
	}
	return true
}

func newServer() (*Server, error) {
	clientCertPath := os.Getenv("CF_INSTANCE_CERT")
	clientKeyPath := os.Getenv("CF_INSTANCE_KEY")

	_, err := os.Stat(clientCertPath)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(clientKeyPath)
	if err != nil {
		return nil, err
	}

	clientCertificate, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, err
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{clientCertificate},
	}

	transport := &http.Transport{TLSClientConfig: tlsConf}
	client := &http.Client{Transport: transport}

	return &Server{client: client}, err
}
