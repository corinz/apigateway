package app

import (
	"encoding/json"
	"io/ioutil"
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
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint(w, r)
		if err := a.MarshalSave(); err != nil {
			errHandler(&w, http.StatusInternalServerError, "ERROR: createAPIEndpoint: Unable to marshal and save app data")
			// TODO revert changes if unable to save
		}
	}
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
	a.apis.AddAPI(api)
	apiPtr := a.apis.GetAPI(api.Name)

	// Create default endpoint and append to new api
	defaultRequest := agw.Request{RequestBody: "", RequestURL: "https://httpbin.org/get", RequestVerb: "GET"}
	defaultEndpoint := agw.APIEndpoint{Name: "default", Description: "default endpoint", ParentName: apiPtr.Name, Request: defaultRequest}
	apiPtr.AppendEndpoint(defaultEndpoint)

	json.NewEncoder(w).Encode(apiPtr)
}

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func (a *app) executeAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	// Execute endpoint
	resp, err := a.apis.GetAPI(apiName).GetAPIEndpoint(endpoint).Execute()
	if err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: executeAPIEndpoint:"+err.Error())
		return
	}
	defer resp.Body.Close() // Response body not closed by Execute()

	if resp.StatusCode == http.StatusOK { // if OK, get & write resp
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errHandler(&w, http.StatusInternalServerError, "ERROR: executeAPIEndpoint:"+err.Error())
			return
		}
		bodyString := string(bodyBytes)
		json.NewEncoder(w).Encode(bodyString)
	} else {
		errHandler(&w, resp.StatusCode, "ERROR: executeAPIEndpoint:"+resp.Status)
		return
	}
}

// listAPIs writes json encoded apis struct to the response writer
// ../
func (a *app) listAPIs(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.apis.APIArr); err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPIs: "+err.Error())
	}
}

// listAPI writes json encoded api struct to the response writer
// ../{api}
func (a *app) listAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	if err := json.NewEncoder(w).Encode(a.apis.GetAPI(apiName)); err != nil {
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
	if err := json.NewEncoder(w).Encode(ep); err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPIEndpoints: "+err.Error())
	}
}
