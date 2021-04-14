package apigateway

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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
	if err := json.Unmarshal(body, &a); err != nil {
		log.Printf(err.Error())
		return a, err
	}
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
	var aep APIEndpoint
	if json.Valid(body) == false {
		err := errors.New("ERROR: unmarshalAPIEndpoint: JSON Invalid")
		log.Printf(err.Error())
		return aep, err
	}
	if err := json.Unmarshal(body, &aep); err != nil {
		log.Printf(err.Error())
		return aep, err
	}
	if aep.Name == "" || aep.Request.RequestURL == "" || aep.Request.RequestVerb == "" {
		err := errors.New("ERROR: unmarshalAPIEndpoint: Required parm missing")
		log.Printf(err.Error())
		return aep, err
	}
	return aep, nil
}
