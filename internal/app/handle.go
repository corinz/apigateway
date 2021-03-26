package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	agw "github.com/corinz/apigateway/pkg/apigateway"

	"github.com/gorilla/mux"
)

func errHandler(w *http.ResponseWriter, errCode int, errStr string) {
	log.Printf(errStr)
	http.Error(*w, errStr, errCode)
}

// record is a handlefunc decorator that marshals and saves the APIs struct to a file
func (a *app) record(endpoint func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			endpoint(w, r)
			err := a.MarshalSave()
			if err != nil {
				errHandler(&w, http.StatusInternalServerError, "ERROR: createAPIEndpoint: Unable to marshal and save app data")
				// TODO revert changes if unable to save
			}
		})
}

// createAPIEndpoint creates an an apiEndpoint from POST data and appends to the api named in the path
// ../{api}/{endpoint}
func (a *app) createAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]

	// error and return if API does not exist
	apiPtr := a.apis.GetAPI(apiName)
	if apiPtr == nil {
		errHandler(&w, http.StatusNotFound, "ERROR: createAPIEndpoint: Requested API object does not exist")
		return
	}

	// init endpoint, error and return if endpoint is invalid or exists
	apiEP, err := agw.UnmarshalAPIEndpoint(r)
	if err != nil {
		errHandler(&w, http.StatusBadRequest, "ERROR: createAPIEndpoint: "+err.Error())
		return
	}
	apiEP.ParentName = apiPtr.Name
	if a.apis.Exists(apiEP) {
		errHandler(&w, 409, "ERROR: createAPIEndpoint: Requested API Endpoint object exists")
		return
	}

	apiPtr.AppendEndpoint(apiEP)
	a.newHandleFunc(apiPtr, apiEP.Name)

	json.NewEncoder(w).Encode(apiPtr.GetAPIEndpoint(apiEP.Name))
}

// createAPI
// ../{api}
func (a *app) createAPI(w http.ResponseWriter, r *http.Request) {
	api, err := agw.UnmarshalAPI(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if a.apis.Exists(api) {
		errHandler(&w, http.StatusConflict, "ERROR: createAPI: Requested API object exists")
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
		ParentName:  apiPtr.Name,
	}
	apiPtr.AppendEndpoint(root)

	// Add router & handle
	a.addSubRouter(apiPtr)
	a.newHandleFunc(apiPtr, "")

	json.NewEncoder(w).Encode(apiPtr)
}

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func (a *app) executeAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	err := a.apis.GetAPI(apiName).GetAPIEndpoint(endpoint).Execute()
	if err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: executeAPIEndpoint:"+err.Error())
		return
	}
}

// listAPIs writes json encoded apis struct to the response writer
// ../
func (a *app) listAPIs(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(a.apis.APIArr)
	if err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPIs: "+err.Error())
	}
}

// listAPI writes json encoded api struct to the response writer
// ../{api}
func (a *app) listAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	err := json.NewEncoder(w).Encode(a.apis.GetAPI(apiName))
	if err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPI: "+err.Error())
	}
}

// listAPIEndpoints writes json encoded apiEndpoint struct to the response writer
// ../{api}/{endpoint}
func (a *app) listAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep := a.apis.GetAPI(apiName).GetAPIEndpoint(epName)
	err := json.NewEncoder(w).Encode(ep) //TODO encoding doesnt work for this struct
	if err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPIEndpoints: "+err.Error())
	}
}

// generic is a placeholder method
func generic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the generic method executing...")
}
