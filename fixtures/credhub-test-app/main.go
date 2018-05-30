package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const credhubBaseURL = "https://credhub.service.cf.internal:8844"

func main() {

	s, err := newServer()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/create", s.Create)
	http.HandleFunc("/list", s.List)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

type Server struct {
	client  *http.Client
	counter int
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	s.counter++
	resp, err := s.client.Post(fmt.Sprintf("%s/api/v1/data", credhubBaseURL), "application/json",
		strings.NewReader(fmt.Sprintf(`{"name": "/%s/%d", "type": "password"}`, os.Getenv("CF_INSTANCE_GUID"), s.counter)),
	)
	if ok := handleBadResponses(w, resp, err); !ok {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	resp, err := s.client.Get(fmt.Sprintf("%s/api/v1/data?path=%s", credhubBaseURL, os.Getenv("CF_INSTANCE_GUID")))
	if ok := handleBadResponses(w, resp, err); !ok {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Encountered error reading response body from credhub: [%s]", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, string(body))
}

func handleBadResponses(w http.ResponseWriter, resp *http.Response, err error) bool {
	if err != nil {
		fmt.Fprintf(w, "Encountered error from credhub: [%s]", err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusBadRequest {
		fmt.Fprintf(w, "Encountered bad response code from credhub: %d", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
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
