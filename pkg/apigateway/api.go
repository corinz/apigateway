package apigateway

import (
	"log"
	"net/http"
	"strings"
)

// APIs represents a slice of type: API
type APIs struct {
	APIArr []API
}

// API represents a slice of type: APIEndpoint
type API struct {
	Name   string `json:"Name"`
	APIEPs []APIEndpoint
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

// Exists checks to see if an interface of type api or apiEndpoint exist
func (apis *APIs) Exists(thing interface{}) bool {
	switch thing.(type) {
	case API:
		a := thing.(API)
		if apis.GetAPI(a.Name) != nil { // if API found
			return true
		}
	case APIEndpoint:
		aep := thing.(APIEndpoint)
		api := apis.GetAPI(aep.ParentName)
		if api.GetAPIEndpoint(aep.Name) != nil { // if API Endpoint found
			return true
		}
	}
	return false
}

// addAPI appends to api slice
func (apis *APIs) AddAPI(a API) {
	apis.APIArr = append(apis.APIArr, a)
}

// getAPI accepts name argument and returns pointer to api
func (apis *APIs) GetAPI(name string) (*API) {
	for i, api := range apis.APIArr {
		if api.Name == name {
			return &apis.APIArr[i]
		}
	}
	return nil
}

// getAPI accepts name argument and returns pointer to api and index
func (apis *APIs) GetAPIIndex(name string) (*API, int) {
	for i, api := range apis.APIArr {
		if api.Name == name {
			return &apis.APIArr[i], i
		}
	}
	return nil, -1
}

// appendEndpoint appends an endpoint to the apiEPs slice
func (api *API) AppendEndpoint(aep APIEndpoint) {
	api.APIEPs = append(api.APIEPs, aep)
}

// getAPIEndpoint returns a pointer to the endpoint struct or nil if not found
func (api *API) GetAPIEndpoint(apiEPName string) (*APIEndpoint) {
	for _, apiEP := range api.APIEPs {
		if apiEP.Name == apiEPName {
			return &apiEP
		}
	}
	return nil
}

// getAPIEndpoint returns a pointer to the endpoint struct or nil if not found
func (api *API) GetAPIEndpointIndex(apiEPName string) (*APIEndpoint, int) {
	for i, apiEP := range api.APIEPs {
		if apiEP.Name == apiEPName {
			return &apiEP, i
		}
	}
	return nil, -1
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
