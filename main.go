package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

type apiRouter struct {
	r    *mux.Router
	apis []api
}

// api is a struct representing APIEndpoints
type api struct {
	Name   string `json:"Name"`
	Route  string `json:"Route"`
	apiEPs []apiEndpoint
	router *mux.Router
}

// apiEndpoint is a struct representing a single API Endpoint with a route and http verb
type apiEndpoint struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HTTPVerb    string `json:"HTTPVerb"`
	Route       string `json:"Route"`
	Command     string `json:"Command"`
	UID         int
}

// createAPIEndpoint creates an an apiEndpoint from POST data and appends to the api named in the path
// ../{api}/{endpoint}
func (ar *apiRouter) createAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	apiEP := unmarshalAPIEndpoint(r)

	ar.appendEndpoint(apiName, apiEP)
	ar.newHandleFunc(apiName)
	json.NewEncoder(w).Encode(apiEP)
}

// createAPI
// ../{api}
func (ar *apiRouter) createAPI(w http.ResponseWriter, r *http.Request) {
	api := unmarshalAPI(r) //TODO Test to see if API already exists and validate the Route is unique/valid

	// Create root endpoint and append to new api
	apiRootEP := apiEndpoint{
		Name:        "default",
		Description: "default endpoint",
		HTTPVerb:    "GET",
		Route:       "/",
		UID:         0,
		Command:     "whoami",
	}
	api.apiEPs = append(api.apiEPs, apiRootEP) // TODO write a setter method for this
	ar.apis = append(ar.apis, api)             // TODO write a setter method for this

	// create new subrouter
	ar.newAPIRouter(api.Name)

	// add new HandleFunc to subrouter
	ar.newHandleFunc(api.Name) // TODO Test Route var syntax and if it has been used already
	json.NewEncoder(w).Encode(&api)
}

// newAPIRouter creates a subrouter from the parent router
func (ar *apiRouter) newAPIRouter(name string) {
	api := ar.getAPI(name)
	api.router = ar.r.PathPrefix("/" + name).Subrouter() // "/{apiName}/"
}

// newHandleFunc creates a new route on the subrouter
func (ar *apiRouter) newHandleFunc(apiName string) {
	api := ar.getAPI(apiName)
	if api == nil {
		log.Printf("API %v not found", apiName)
	}
	api.router.HandleFunc("/"+apiName, generic) // "/{apiName}/{aepName}/"
}

// appendEndpoint appends an endpoint to the apiEPs slice
func (ar *apiRouter) appendEndpoint(apiName string, apiEP apiEndpoint) {
	api := ar.getAPI(apiName)
	api.apiEPs = append(api.apiEPs, apiEP)
}

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func (ar *apiRouter) executeAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	apiEP := ar.getAPIEndpoint(apiName, endpoint)
	err := apiEP.execute() // TODO handle error, return error code/message
	json.NewEncoder(w).Encode(err)
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

// getAPIEndpoint accepts api & apiEndpoint name and returns a pointer to the apiEndpoint
func (ar *apiRouter) getAPIEndpoint(apiName string, apiEPName string) *apiEndpoint {
	for _, api := range ar.apis {
		if api.Name == apiName {
			for _, apiEP := range api.apiEPs {
				if apiEP.Name == apiEPName {
					return &apiEP
				}
			}
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
	api := ar.getAPI(apiName)
	json.NewEncoder(w).Encode(api)
}

// listAPIEndpoints writes json encoded apiEndpoint struct to the response writer
// ../{api}/{endpoint}
func (ar *apiRouter) listAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep := ar.getAPIEndpoint(apiName, epName)
	json.NewEncoder(w).Encode(ep)
}

// execute executes the command found in the apiEndpoint.Command struct-field
func (aep *apiEndpoint) execute() error {
	cmd := exec.Command("sleep", "30") // TODO this is hardcoded for now
	log.Printf("Running command...")
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
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
func unmarshalAPI(r *http.Request) api {
	//TODO error handle and payload validation
	body, _ := ioutil.ReadAll(r.Body)
	var a api
	json.Unmarshal(body, &a)

	return a
}

// unmarshalAPIEndpoint accepts http request and returns unmarshalled apiEndpoint struct
func unmarshalAPIEndpoint(r *http.Request) apiEndpoint {
	//TODO error handle and payload validation
	body, _ := ioutil.ReadAll(r.Body)
	var apiEP apiEndpoint
	json.Unmarshal(body, &apiEP)

	return apiEP
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
