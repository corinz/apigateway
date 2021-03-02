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

// apiEndpoint is a struct representing a single API Endpoint with a route and http verb
type apiEndpoint struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HTTPVerb    string `json:"HTTPVerb"`
	Route       string `json:"Route"`
	Command     string `json:"Command"`
	UID         int
}

// api is a struct representing APIEndpoints
type api struct {
	Name   string `json:"Name"`
	Route  string `json:"Route"`
	apiEPs []apiEndpoint
	router *mux.Router
}

// TODO struct that wraps globals, instantiate it in main()

// APIs is a global var representing all of the applictions api structs
var apis []api

// mainRouter is the parent router of all application routers
var mainRouter = mux.NewRouter().StrictSlash(true)

// createAPIEndpoint creates an an apiEndpoint from POST data and appends to the api named in the path
// ../{api}/{endpoint}
func createAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	api := getAPI(apiName) //TODO return error if not found
	apiEP := unmarshalAPIEndpoint(r)

	appendAPIEndpoint(apiName, apiEP)
	newHandleFunc(api, apiEP.Name)
	json.NewEncoder(w).Encode(apiEP)
}

// appendAPIEndpoint adds an endpoint to the apiEPs slice
func appendAPIEndpoint(apiName string, apiEP apiEndpoint) {
	api := getAPI(apiName)
	api.apiEPs = append(api.apiEPs, apiEP)
}

// createAPI
// ../{api}
func createAPI(w http.ResponseWriter, r *http.Request) {
	api := unmarshalAPI(r) //TODO Test to see if API already exists and validate the Route is unique/valid

	// Create root endpoint and append to new api
	apiRootEP := apiEndpoint{
		Name:        api.Name + "/rootEP",
		Description: api.Name + "/rootEP",
		HTTPVerb:    "GET",
		Route:       "/",
		UID:         0,
		Command:     "whoami",
	}
	api.apiEPs = append(api.apiEPs, apiRootEP) // TODO write a setter method for this
	apis = append(apis, api)                   // TODO write a setter method for this
	apiPtr := getAPI(api.Name)

	// create new subrouter
	newAPIRouter(apiPtr)

	// add new HandleFunc to subrouter
	newHandleFunc(apiPtr, api.apiEPs[0].Name) // TODO Test Route var syntax and if it has been used already
	json.NewEncoder(w).Encode(&api)
}

// listAPIs writes json encoded apis struct to the response writer
// ../
func listAPIs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(apis)
}

// listAPI writes json encoded api struct to the response writer
// ../{api}
func listAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	api := getAPI(apiName)
	json.NewEncoder(w).Encode(api)
}

// listAPIEndpoints writes json encoded apiEndpoint struct to the response writer
// ../{api}/{endpoint}
func listAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep := getAPIEndpoint(apiName, epName)
	json.NewEncoder(w).Encode(ep)
}

// newAPIRouter creates a subrouter from the parent router
func newAPIRouter(api *api) {
	api.router = mainRouter.PathPrefix("/" + api.Name).Subrouter() // "/{apiName}/"
}

// newHandleFunc creates a new route on the subrouter
func newHandleFunc(a *api, aepName string) {
	a.router.HandleFunc("/"+aepName, generic) // "/{apiName}/{aepName}/"
}

// generic is a placeholder method
func generic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the generic method executing...")
}

// execute executes the command found in the apiEndpoint.Command struct-field
func execute(apiEP *apiEndpoint) error {
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

// executeAPIEndpoint locates the apiEndpoint struct and calls execute()
func executeAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	apiEP := getAPIEndpoint(apiName, endpoint)
	err := execute(apiEP) // TODO handle error, return error code/message
	json.NewEncoder(w).Encode(err)
}

// getAPI accepts name argument and returns pointer to api
func getAPI(name string) *api {
	for i, api := range apis {
		if api.Name == name {
			return &apis[i]
		}
	}
	return nil
}

// getAPIEndpoint accepts api & apiEndpoint name and returns a pointer to the apiEndpoint
func getAPIEndpoint(apiName string, apiEPName string) *apiEndpoint {
	for _, api := range apis {
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

func startup() {
	// GETs
	mainRouter.HandleFunc("/", listAPIs).Methods("GET")
	mainRouter.HandleFunc("/{api}", listAPI).Methods("GET")
	mainRouter.HandleFunc("/{api}/{endpoint}", listAPIEndpoints).Methods("GET")

	// POSTs
	mainRouter.HandleFunc("/", createAPI).Methods("POST")
	mainRouter.HandleFunc("/{api}", createAPIEndpoint).Methods("POST")
	mainRouter.HandleFunc("/{api}/{endpoint}", executeAPIEndpoint).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", mainRouter))
}

func main() {
	startup()
}
