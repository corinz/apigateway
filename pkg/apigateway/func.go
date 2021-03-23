package apigateway

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	//"github.com/gorilla/mux"
)

// unmarshalAPI accepts http request and returns unmarshalled api struct
// Checks if json is valid and if 'Name' parm exists
func UnmarshalAPI(r *http.Request) (API, error) {
	//TODO Combine with other unmarshal func
	body, _ := ioutil.ReadAll(r.Body)
	var a API
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
func UnmarshalAPIEndpoint(r *http.Request) (APIEndpoint, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var apiEP APIEndpoint
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
