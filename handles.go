package APIGateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Startup() {
	apiGW := newAPIGateway()

	// GETs
	apiGW.r.HandleFunc("/", apiGW.listAPIs).Methods("GET")
	apiGW.r.HandleFunc("/{api}", apiGW.listAPI).Methods("GET")
	apiGW.r.HandleFunc("/{api}/{endpoint}", apiGW.listAPIEndpoints).Methods("GET")

	// POSTs
	apiGW.r.HandleFunc("/", apiGW.createAPI).Methods("POST")
	apiGW.r.HandleFunc("/{api}", apiGW.createAPIEndpoint).Methods("POST")
	apiGW.r.HandleFunc("/{api}/{endpoint}", apiGW.executeAPIEndpoint).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", apiGW.r))
}

// createAPIEndpoint creates an an apiEndpoint from POST data and appends to the api named in the path
// ../{api}/{endpoint}
func (ar *apiRouter) createAPIEndpoint(w http.ResponseWriter, r *http.Request) {
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
func (ar *apiRouter) createAPI(w http.ResponseWriter, r *http.Request) {
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
	root := apiEndpoint{
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

// exists checks to see if an interface of type api or apiEndpoint exist
func (ar *apiRouter) exists(thing interface{}) bool {

	switch thing.(type) {
	case api:
		a := thing.(api)
		if ar.getAPI(a.Name) != nil { // if API found
			return true
		}
	case apiEndpoint:
		aep := thing.(apiEndpoint)
		api := aep.parentPtr
		if api.getAPIEndpoint(aep.Name) != nil { // if API Endpoint found
			return true
		}
	}
	return false
}

// addAPI appends to api slice
func (ar *apiRouter) addAPI(a api) {
	ar.apis = append(ar.apis, a)
}

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func (ar *apiRouter) executeAPIEndpoint(w http.ResponseWriter, r *http.Request) {
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
func (ar *apiRouter) listAPIs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(ar.apis)
}

// listAPI writes json encoded api struct to the response writer
// ../{api}
func (ar *apiRouter) listAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	json.NewEncoder(w).Encode(ar.getAPI(apiName).apiEPs)
}

// listAPIEndpoints writes json encoded apiEndpoint struct to the response writer
// ../{api}/{endpoint}
func (ar *apiRouter) listAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep := ar.getAPI(apiName).getAPIEndpoint(epName)
	json.NewEncoder(w).Encode(ep)
}

// generic is a placeholder method
func generic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the generic method executing...")
}
