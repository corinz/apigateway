package apigateway

import (
	"log"
	"os/exec"
	"strings"
)

// getAPI accepts name argument and returns pointer to api
func (ar *apiRouter) getAPI(name string) *api {
	for i, api := range ar.apis {
		if api.Name == name {
			return &ar.apis[i]
		}
	}
	return nil
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
