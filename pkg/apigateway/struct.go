package apigateway

import "github.com/gorilla/mux"

type apiRouter struct {
	r    *mux.Router
	apis []api
}

// api is a struct representing APIEndpoints
type api struct {
	Name   string `json:"Name"`
	apiEPs []apiEndpoint
	router *mux.Router
}

// apiEndpoint is a struct representing a single API Endpoint with a route and http verb
type apiEndpoint struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HTTPVerb    string `json:"HTTPVerb"`
	Command     string `json:"Command"`
	UID         int
	parentPtr   *api
}
