package apigateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// createAPIEndpoint creates an an apiEndpoint from POST data and appends to the api named in the path
// ../{api}/{endpoint}
func (ar *apiRouter) CreateAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]

	// error and return if API does not exist
	a := ar.getAPI(apiName)
	if a == nil {
		errStr := "ERROR: createAPIEndpoint: Requested API object does not exist"
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusNotFound)
		return
	}

	// init endpoint, error and return if endpoint is invalid or exists
	apiEP, err := unmarshalAPIEndpoint(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	apiEP.parentPtr = a
	if ar.exists(apiEP) {
		errStr := "ERROR: createAPIEndpoint: Requested API Endpoint object exists"
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusConflict)
		return
	}

	a.appendEndpoint(apiEP)
	a.newHandleFunc()
	json.NewEncoder(w).Encode(apiEP)
}

// createAPI
// ../{api}
func (ar *apiRouter) CreateAPI(w http.ResponseWriter, r *http.Request) {
	api, err := unmarshalAPI(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ar.exists(api) {
		errStr := "ERROR: createAPI: Requested API object exists"
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusConflict)
		return
	}

	// Appends api
	ar.addAPI(api)

	// Create root endpoint and append to new api
	root := ApiEndpoint{
		Name:        "default",
		Description: "default endpoint",
		HTTPVerb:    "GET",
		UID:         0,
		Command:     "whoami",
	}

	a := ar.getAPI(api.Name)
	a.appendEndpoint(root)
	a.newAPIRouter(ar)
	a.newHandleFunc()

	json.NewEncoder(w).Encode(&api)
}

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func (ar *apiRouter) ExecuteAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	err := ar.getAPI(apiName).getAPIEndpoint(endpoint).execute()
	if err != nil {
		errStr := "ERROR: executeAPIEndpoint:" + err.Error()
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}
}

// listAPIs writes json encoded apis struct to the response writer
// ../
func (ar *apiRouter) ListAPIs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(ar.apis)
}

// listAPI writes json encoded api struct to the response writer
// ../{api}
func (ar *apiRouter) ListAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	json.NewEncoder(w).Encode(ar.getAPI(apiName).apiEPs)
}

// listAPIEndpoints writes json encoded apiEndpoint struct to the response writer
// ../{api}/{endpoint}
func (ar *apiRouter) ListAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep := ar.getAPI(apiName).getAPIEndpoint(epName)
	json.NewEncoder(w).Encode(ep)
}

// generic is a placeholder method
func Generic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the generic method executing...")
}
