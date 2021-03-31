package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Marshal converts APIs struct to json and saves locally
func (a *app) MarshalSave() error {
	jsonByte, err := json.Marshal(a.apis)
	if err != nil {
		log.Println(err)
		return err
	}
	if err = ioutil.WriteFile("apigateway.json", jsonByte, 0644); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
