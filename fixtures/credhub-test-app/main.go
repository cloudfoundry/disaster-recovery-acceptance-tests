package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/generate"
	"strings"
)

const credhubBaseURL = "https://credhub.service.cf.internal:8844"
const credentialName = "SECRET_PASSWORD"

func main() {
	s, err := newServer()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/create", s.Create)
	http.HandleFunc("/list", s.List)
	http.HandleFunc("/clean", s.Clean)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

type Server struct {
	client  *credhub.CredHub
	counter int
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	s.counter++
	password, err := s.client.GeneratePassword(fmt.Sprintf("/%s/%d",credentialName,s.counter),generate.Password{},credhub.Overwrite)
	if ok := handleBadResponses(w, err); !ok {
		return
	}
	fmt.Println(password)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	results, err := s.client.FindByPartialName("/"+credentialName)
	if ok := handleBadResponses(w, err); !ok {
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered error reading response body from credhub: [%s]", err)
	}

	var response []string
	for _, cred := range results.Credentials {
		response = append(response, "{\"name\": \"" + cred.Name+ "\"}")
	}

	respString := strings.Join(response, ",")


	fmt.Fprintf(w, "{ \"credentials\": ["+string(respString)+"]}")
}

func (s *Server) Clean(w http.ResponseWriter, r *http.Request) {
	results, err := s.client.FindByPartialName("/"+credentialName)

	if ok := handleBadResponses(w, err); !ok {
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered error reading response body from credhub: [%s]", err)
	}

	for _, cred := range results.Credentials {
		err = s.client.Delete(cred.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Encountered error reading response body from credhub: [%s]", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}


func handleBadResponses(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered error from credhub: [%s]", err)
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

	client, err := credhub.New(
		credhubBaseURL,
		credhub.SkipTLSValidation(true),
		credhub.Auth(auth.UaaClientCredentials(os.Getenv("CREDHUB_CLIENT"), os.Getenv("CREDHUB_SECRET"))))

	return &Server{client: client}, err
}
