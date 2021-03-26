package apigateway

import "github.com/gorilla/mux"

//TODO get rid of this struct?
type APIs struct {
	APIArr []API
}

// API is a struct representing APIEndpoints
type API struct {
	Name   string `json:"Name"`
	APIEPs []APIEndpoint
	Router *mux.Router
}

// APIEndpoint is a struct representing a single API Endpoint with a route and http verb
type APIEndpoint struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HTTPVerb    string `json:"HTTPVerb"`
	Command     string `json:"Command"`
	UID         int
	ParentName  string
}
