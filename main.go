package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

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

// getAPI accepts name argument and returns pointer to api
func (ar *apiRouter) getAPI(name string) *api {
	for i, api := range ar.apis {
		if api.Name == name {
			return &ar.apis[i]
		}
	}
	return nil
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

// newAPIRouter creates a subrouter from the parent router
func (a *api) newAPIRouter(ar *apiRouter) {
	a.router = ar.r.PathPrefix("/" + a.Name).Subrouter() // "/{apiName}/"
}

// newHandleFunc creates a new route on the subrouter
func (a *api) newHandleFunc() {
	a.router.HandleFunc("/"+a.Name, generic) // "/{apiName}/{aepName}/"
}

// appendEndpoint appends an endpoint to the apiEPs slice
func (a *api) appendEndpoint(aep apiEndpoint) {
	a.apiEPs = append(a.apiEPs, aep)
}

// getAPIEndpoint returns a pointer to the endpoint struct or nil if not found
func (a *api) getAPIEndpoint(apiEPName string) *apiEndpoint {
	for _, apiEP := range a.apiEPs {
		if apiEP.Name == apiEPName {
			return &apiEP
		}
	}
	return nil
}

// execute executes the command found in the apiEndpoint.Command struct-field
func (aep *apiEndpoint) execute() error {
	s := strings.Split(aep.Command, " ")

	cmd := exec.Command(s[0], s[1:]...)
	log.Printf("Running command: %v", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Printf("Execution complete")
	return nil
}

// generic is a placeholder method
func generic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the generic method executing...")
}

// unmarshalAPI accepts http request and returns unmarshalled api struct
// Checks if json is valid and if 'Name' parm exists
func unmarshalAPI(r *http.Request) (api, error) {
	//TODO Combine with other unmarshal func
	body, _ := ioutil.ReadAll(r.Body)
	var a api
	if json.Valid(body) == false {
		err := errors.New("ERROR: unmarshalAPI: JSON Invalid")
		log.Printf(err.Error())
		return a, err
	}
	json.Unmarshal(body, &a)

	if a.Name == "" {
		err := errors.New("ERROR: unmarshalAPI: Required parm missing")
		log.Printf(err.Error())
		return a, err
	}
	return a, nil
}

// unmarshalAPIEndpoint accepts http request and returns unmarshalled apiEndpoint struct
func unmarshalAPIEndpoint(r *http.Request) (apiEndpoint, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var apiEP apiEndpoint
	if json.Valid(body) == false {
		err := errors.New("ERROR: unmarshalAPIEndpoint: JSON Invalid")
		log.Printf(err.Error())
		return apiEP, err
	}
	json.Unmarshal(body, &apiEP)

	if apiEP.Name == "" {
		err := errors.New("ERROR: unmarshalAPIEndpoint: Required parm missing")
		log.Printf(err.Error())
		return apiEP, err
	}

	return apiEP, nil
}

// NewAPIGateway inits a new apiRouter struct
func newAPIGateway() *apiRouter {
	r := mux.NewRouter().StrictSlash(true)
	return &apiRouter{r: r}
}

func main() {
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
