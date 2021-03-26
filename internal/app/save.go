package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Marshal converts APIs struct to json and saves locally
func (a *app) MarshalSave() error {
	json, err := json.Marshal(a.apis)
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile("apigateway.json", json, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
