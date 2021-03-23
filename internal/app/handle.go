package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	agw "github.com/corinz/apigateway/pkg/apigateway"

	"github.com/gorilla/mux"
)

// createAPIEndpoint creates an an apiEndpoint from POST data and appends to the api named in the path
// ../{api}/{endpoint}
func (a *app) CreateAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]

	// error and return if API does not exist
	apiPtr := a.apis.GetAPI(apiName)
	if apiPtr == nil {
		errStr := "ERROR: createAPIEndpoint: Requested API object does not exist"
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusNotFound)
		return
	}

	// init endpoint, error and return if endpoint is invalid or exists
	apiEP, err := agw.UnmarshalAPIEndpoint(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	apiEP.ParentPtr = apiPtr
	if a.apis.Exists(apiEP) {
		errStr := "ERROR: createAPIEndpoint: Requested API Endpoint object exists"
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusConflict)
		return
	}

	apiPtr.AppendEndpoint(apiEP)
	a.newHandleFunc(apiPtr)
	json.NewEncoder(w).Encode(apiPtr.GetAPIEndpoint(apiEP.Name)) //TODO validate this new output var
}

// createAPI
// ../{api}
func (a *app) CreateAPI(w http.ResponseWriter, r *http.Request) {
	api, err := agw.UnmarshalAPI(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if a.apis.Exists(api) {
		errStr := "ERROR: createAPI: Requested API object exists"
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusConflict)
		return
	}

	// Appends api
	a.apis.AddAPI(api)

	apiPtr := a.apis.GetAPI(api.Name)

	// Create root endpoint and append to new api
	root := agw.APIEndpoint{
		Name:        "default",
		Description: "default endpoint",
		HTTPVerb:    "GET",
		UID:         0,
		Command:     "whoami",
		// ParentPtr:   apiPtr, // TODO this prevents object from being necoded and returned
	}
	apiPtr.AppendEndpoint(root)

	// Add router & handle
	a.addSubRouter(apiPtr)
	a.newHandleFunc(apiPtr)

	json.NewEncoder(w).Encode(apiPtr)
}

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func (a *app) ExecuteAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	err := a.apis.GetAPI(apiName).GetAPIEndpoint(endpoint).Execute()
	if err != nil {
		errStr := "ERROR: executeAPIEndpoint:" + err.Error()
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}
}

// listAPIs writes json encoded apis struct to the response writer
// ../
func (a *app) ListAPIs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(a.apis.APIArr)
}

// listAPI writes json encoded api struct to the response writer
// ../{api}
func (a *app) ListAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	json.NewEncoder(w).Encode(a.apis.GetAPI(apiName).APIEPs)
}

// listAPIEndpoints writes json encoded apiEndpoint struct to the response writer
// ../{api}/{endpoint}
func (a *app) ListAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep := a.apis.GetAPI(apiName).GetAPIEndpoint(epName)
	fmt.Println("ep:", ep)
	json.NewEncoder(w).Encode(ep) //TODO encoding doesnt work for this struct
}

// generic is a placeholder method
func Generic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the generic method executing...")
}
