package apigateway

import (
	"log"
	"os/exec"
	"strings"
)

// exists checks to see if an interface of type api or apiEndpoint exist
func (apis *APIs) Exists(thing interface{}) bool {

	switch thing.(type) {
	case API:
		a := thing.(API)
		if apis.GetAPI(a.Name) != nil { // if API found
			return true
		}
	case APIEndpoint:
		aep := thing.(APIEndpoint)
		api := aep.ParentPtr
		if api.GetAPIEndpoint(aep.Name) != nil { // if API Endpoint found
			return true
		}
	}
	return false
}

// addAPI appends to api slice
func (apis *APIs) AddAPI(a API) {
	apis.APIArr = append(apis.APIArr, a)
}

// getAPI accepts name argument and returns pointer to api
func (apis *APIs) GetAPI(name string) *API {
	for i, api := range apis.APIArr {
		if api.Name == name {
			return &apis.APIArr[i]
		}
	}
	return nil
}

// appendEndpoint appends an endpoint to the apiEPs slice
func (api *API) AppendEndpoint(aep APIEndpoint) {
	api.APIEPs = append(api.APIEPs, aep)
}

// getAPIEndpoint returns a pointer to the endpoint struct or nil if not found
func (api *API) GetAPIEndpoint(apiEPName string) *APIEndpoint {
	for _, apiEP := range api.APIEPs {
		if apiEP.Name == apiEPName {
			return &apiEP
		}
	}
	return nil
}

// execute executes the command found in the apiEndpoint.Command struct-field
func (aep *APIEndpoint) Execute() error {
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
