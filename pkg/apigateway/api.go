package apigateway

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

// APIs represents a slice of type: API
type APIs struct {
	APIMap map[string]API
	sync.RWMutex
}

// API represents a slice of type: APIEndpoint
type API struct {
	Name     string `json:"Name"`
	APIEPMap map[string]APIEndpoint
}

// APIEndpoint represents a single API Endpoint and is populated by the API user
type APIEndpoint struct {
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Request     Request `json:"Request"`
	ParentName  string
}

// Request represents the users outgoing request and is populated by the API user
type Request struct {
	RequestBody string `json:"RequestBody"`
	RequestURL  string `json:"RequestURL"`
	RequestVerb string `json:"RequestVerb"`
}

// Execute builds the request and executes it
// Note: response body is not closed in this method
//   close response body in caller method with resp.Body.Close()
func (aep *APIEndpoint) Execute() (*http.Response, error) {
	// Build request
	r := aep.Request
	req, err := http.NewRequest(r.RequestVerb, r.RequestURL, strings.NewReader(r.RequestBody))
	log.Printf("Executing endpoint '../%+v/%+v' with parameters: '%+v'.\n", aep.ParentName, aep.Name, r)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
