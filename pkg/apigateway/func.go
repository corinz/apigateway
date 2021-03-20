package apigateway

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NewAPIGateway inits a new apiRouter struct
func NewAPIGateway() *apiRouter {
	r := mux.NewRouter().StrictSlash(true)
	return &apiRouter{R: r}
}

// unmarshalAPI accepts http request and returns unmarshalled api struct
// Checks if json is valid and if 'Name' parm exists
func unmarshalAPI(r *http.Request) (Api, error) {
	//TODO Combine with other unmarshal func
	body, _ := ioutil.ReadAll(r.Body)
	var a Api
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
func unmarshalAPIEndpoint(r *http.Request) (ApiEndpoint, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var apiEP ApiEndpoint
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
