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

	a := ar.getAPI(apiName)
	a.appendEndpoint(apiEP)
	a.newHandleFunc()
	json.NewEncoder(w).Encode(apiEP)
}

// createAPI
// ../{api}
func (ar *apiRouter) createAPI(w http.ResponseWriter, r *http.Request) {
	api := unmarshalAPI(r) //TODO Test to see if API already exists and validate the Route is unique/valid

	// Appends api
	ar.addAPI(api)

	// Create root endpoint and append to new api
	root := apiEndpoint{
		Name:        "default",
		Description: "default endpoint",
		HTTPVerb:    "GET",
		Route:       "/",
		UID:         0,
		Command:     "whoami",
	}

	a := ar.getAPI(api.Name)
	a.appendEndpoint(root)
	a.newAPIRouter(ar)
	a.newHandleFunc()

	json.NewEncoder(w).Encode(&api)
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
	json.NewEncoder(w).Encode(ar.getAPI(apiName))
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
