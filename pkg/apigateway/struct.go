package apigateway

import "github.com/gorilla/mux"

type apiRouter struct {
	R    *mux.Router
	apis []Api
}

// Api is a struct representing APIEndpoints
type Api struct {
	Name   string `json:"Name"`
	apiEPs []ApiEndpoint
	router *mux.Router
}

// ApiEndpoint is a struct representing a single API Endpoint with a route and http verb
type ApiEndpoint struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HTTPVerb    string `json:"HTTPVerb"`
	Command     string `json:"Command"`
	UID         int
	parentPtr   *Api
}
