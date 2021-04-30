package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	agw "github.com/corinz/apigateway/pkg/apigateway"

	"github.com/gorilla/mux"
)

// TODO feature: auth decorator to authorize creation/deletion of endpoints/apis
func errHandler(w *http.ResponseWriter, errCode int, errStr string) {
	log.Printf(errStr)
	http.Error(*w, errStr, errCode)
}

// record is a handlefunc decorator that marshals and saves the APIs struct to a file
func (a *app) record(endpoint func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint(w, r)
		if err := a.marshalSave(); err != nil {
			errHandler(&w, http.StatusInternalServerError, "ERROR: createAPIEndpoint: Unable to marshal and save app data")
			// TODO revert changes if unable to save
		}
	}
}

// createAPIEndpoint creates an an apiEndpoint from POST data and appends to the api named in the path
// ../{api}/{endpoint}
func (a *app) createAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	a.apis.Lock()
	defer a.apis.Unlock()

	vars := mux.Vars(r)
	apiName := vars["api"]

	api, exists := a.apis.APIMap[apiName]
	if !exists {
		errHandler(&w, http.StatusNotFound, "ERROR: createAPIEndpoint: Requested API does not exist")
		return
	}

	// init endpoint, error and return if endpoint is invalid or exists
	apiEP, err := agw.UnmarshalAPIEndpoint(r)
	if err != nil {
		errHandler(&w, http.StatusBadRequest, "ERROR: createAPIEndpoint: "+err.Error())
		return
	}
	apiEP.ParentName = api.Name
	if _, exists := a.apis.APIMap[api.Name].APIEPMap[apiEP.Name]; exists {
		errHandler(&w, 409, "ERROR: createAPIEndpoint: Requested API Endpoint object exists")
		return
	}

	a.apis.APIMap[api.Name].APIEPMap[apiEP.Name] = apiEP
	json.NewEncoder(w).Encode(a.apis.APIMap[api.Name].APIEPMap[apiEP.Name])
}

// createAPI
// ../{api}
func (a *app) createAPI(w http.ResponseWriter, r *http.Request) {
	a.apis.Lock()
	defer a.apis.Unlock()

	api, err := agw.UnmarshalAPI(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := a.apis.APIMap[api.Name]; exists {
		errHandler(&w, http.StatusConflict, "ERROR: createAPI: Requested API object exists")
		return
	}
	if api.APIEPMap != nil { // api.Name only valid request parm in this method
		errHandler(&w, http.StatusBadRequest, "ERROR: createAPI: 'Name' is the only valid field for /api request")
		return
	}

	api.APIEPMap = make(map[string]agw.APIEndpoint)
	a.apis.APIMap[api.Name] = api

	// Create default endpoint and append to new api
	defaultRequest := agw.Request{RequestBody: "", RequestURL: "https://httpbin.org/get", RequestVerb: "GET"}
	defaultEndpoint := agw.APIEndpoint{Name: "default", Description: "default endpoint", ParentName: api.Name, Request: defaultRequest}
	a.apis.APIMap[api.Name].APIEPMap[defaultEndpoint.Name] = defaultEndpoint

	json.NewEncoder(w).Encode(a.apis.APIMap[api.Name])
}

// executeAPI will execute every endpoint in the endpoint slice
func (a *app) executeAPI(w http.ResponseWriter, r *http.Request) {
	a.apis.RLock()
	defer a.apis.RUnlock()

	vars := mux.Vars(r)
	apiName := vars["api"]

	api, exists := a.apis.APIMap[apiName]
	if !exists {
		errHandler(&w, http.StatusNotFound, "ERROR: executeAPI: API not found")
		return
	}
	for _, ep := range api.APIEPMap {
		// Execute endpoint
		resp, err := ep.Execute()
		if err != nil {
			errHandler(&w, http.StatusInternalServerError, "ERROR: executeAPI:"+err.Error())
			return
		}
		writeResp(&w, resp)
		resp.Body.Close() // Response body not closed by Execute()
	}
}

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func (a *app) executeAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	a.apis.RLock()
	defer a.apis.RUnlock()

	vars := mux.Vars(r)
	apiName := vars["api"]
	ep := vars["endpoint"]

	// Get API
	api, exists := a.apis.APIMap[apiName]
	if !exists {
		errHandler(&w, http.StatusNotFound, "ERROR: executeAPIEndpoint: API not found")
		return
	}

	// Get Endpoint
	apiEP, exists := a.apis.APIMap[api.Name].APIEPMap[ep]
	if !exists {
		errHandler(&w, http.StatusNotFound, "ERROR: executeAPIEndpoint: API Endpoint not found")
		return
	}

	// Execute endpoint
	resp, err := apiEP.Execute()
	if err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: executeAPIEndpoint:"+err.Error())
		return
	}
	defer resp.Body.Close() // Response body not closed by Execute()

	writeResp(&w, resp)
}

// writeResp sends response body to response writer
func writeResp(w *http.ResponseWriter, resp *http.Response) {
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errHandler(w, http.StatusInternalServerError, "ERROR: executeAPIEndpoint:"+err.Error())
			return
		}
		bodyString := string(bodyBytes)
		json.NewEncoder(*w).Encode(bodyString)
	} else {
		errHandler(w, resp.StatusCode, "ERROR: executeAPIEndpoint:"+resp.Status)
		return
	}
}

// listAPIs writes json encoded apis struct to the response writer
// ../
func (a *app) listAPIs(w http.ResponseWriter, r *http.Request) {
	a.apis.RLock()
	defer a.apis.RUnlock()

	if err := json.NewEncoder(w).Encode(a.apis.APIMap); err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPIs: "+err.Error())
	}
}

// listAPI writes json encoded api struct to the response writer
// ../{api}
func (a *app) listAPI(w http.ResponseWriter, r *http.Request) {
	a.apis.RLock()
	defer a.apis.RUnlock()

	vars := mux.Vars(r)
	apiName := vars["api"]

	api, exists := a.apis.APIMap[apiName]
	if !exists {
		errHandler(&w, http.StatusNotFound, "ERROR: listAPI: "+"API does not exist")
		return
	}

	if err := json.NewEncoder(w).Encode(api); err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPI: "+err.Error())
	}
}

// listAPIEndpoints writes json encoded apiEndpoint struct to the response writer
// ../{api}/{endpoint}
func (a *app) listAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	a.apis.RLock()
	defer a.apis.RUnlock()

	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep, exists := a.apis.APIMap[apiName].APIEPMap[epName]
	if !exists {
		errHandler(&w, http.StatusNotFound, "ERROR: listAPIEndpoints: "+"API or Endpoint does not exist")
		return
	}
	if err := json.NewEncoder(w).Encode(ep); err != nil {
		errHandler(&w, http.StatusInternalServerError, "ERROR: listAPIEndpoints: "+err.Error())
	}
}

// delete gets the index of the named API/Endpoint and deletes it
func (a *app) delete(w http.ResponseWriter, r *http.Request) {
	a.apis.Lock()
	defer a.apis.Unlock()

	vars := mux.Vars(r)
	apiName := vars["api"]
	ep := vars["endpoint"]

	api, exists := a.apis.APIMap[apiName]
	if !exists {
		errHandler(&w, http.StatusNotFound, "ERROR: delete: Requested API object does not exist")
		return
	}

	if ep == "" { // delete API
		delete(a.apis.APIMap, api.Name)
		return
	} else { // delete API Endpoint
		apiEP, exists := api.APIEPMap[ep]
		if !exists || apiEP.Name == "default" {
			errHandler(&w, http.StatusNotFound, "ERROR: delete: Requested API Endpoint object does not exist or is not available for deletion")
			return
		}
		delete(a.apis.APIMap[apiName].APIEPMap, ep)
		log.Printf("Deleted API Endpoint: '../%v/%v'\n", api.Name, apiEP.Name)
	}
	errHandler(&w, http.StatusNotFound, "ERROR: deleteAPI: Unable to delete API")
}
