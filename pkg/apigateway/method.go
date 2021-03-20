package apigateway

import (
	"log"
	"os/exec"
	"strings"
)

// exists checks to see if an interface of type api or apiEndpoint exist
func (ar *apiRouter) exists(thing interface{}) bool {

	switch thing.(type) {
	case Api:
		a := thing.(Api)
		if ar.getAPI(a.Name) != nil { // if API found
			return true
		}
	case ApiEndpoint:
		aep := thing.(ApiEndpoint)
		api := aep.parentPtr
		if api.getAPIEndpoint(aep.Name) != nil { // if API Endpoint found
			return true
		}
	}
	return false
}

// addAPI appends to api slice
func (ar *apiRouter) addAPI(a Api) {
	ar.apis = append(ar.apis, a)
}

// getAPI accepts name argument and returns pointer to api
func (ar *apiRouter) getAPI(name string) *Api {
	for i, api := range ar.apis {
		if api.Name == name {
			return &ar.apis[i]
		}
	}
	return nil
}

// newAPIRouter creates a subrouter from the parent router
func (a *Api) newAPIRouter(ar *apiRouter) {
	a.router = ar.R.PathPrefix("/" + a.Name).Subrouter() // "/{apiName}/"
}

// newHandleFunc creates a new route on the subrouter
func (a *Api) newHandleFunc() {
	a.router.HandleFunc("/"+a.Name, Generic) // "/{apiName}/{aepName}/"
}

// appendEndpoint appends an endpoint to the apiEPs slice
func (a *Api) appendEndpoint(aep ApiEndpoint) {
	a.apiEPs = append(a.apiEPs, aep)
}

// getAPIEndpoint returns a pointer to the endpoint struct or nil if not found
func (a *Api) getAPIEndpoint(apiEPName string) *ApiEndpoint {
	for _, apiEP := range a.apiEPs {
		if apiEP.Name == apiEPName {
			return &apiEP
		}
	}
	return nil
}

// execute executes the command found in the apiEndpoint.Command struct-field
func (aep *ApiEndpoint) execute() error {
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
