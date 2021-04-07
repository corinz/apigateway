package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Marshal converts APIs struct to json and saves locally
func (a *app) MarshalSave() error {
	a.apis.RLock()
	defer a.apis.RUnlock()

	var err error
	raw, err := json.Marshal(a.apis)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = ioutil.WriteFile("apigateway.json", raw, 0600); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
