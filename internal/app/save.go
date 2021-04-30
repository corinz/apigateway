package app

import (
	"encoding/json"
	"github.com/corinz/apigateway/pkg/apigateway"
	"io/ioutil"
	"log"
)

// Marshal converts APIs struct to json and saves locally
func (a *app) marshalSave() error {
	a.apis.RLock()
	defer a.apis.RUnlock()

	var err error
	raw, err := json.Marshal(a.apis.APIMap)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = ioutil.WriteFile("cmd/server/apigateway.json", raw, 0600); err != nil {
		log.Println("Save error: " + err.Error())
		return err
	}
	return nil
}

// load reads the saved json file and unmarshals to the struct
func (a *app) load() error {
	a.apis.Lock()
	defer a.apis.Unlock()

	var apis apigateway.APIs
	var err error
	var raw []byte

	if raw, err = ioutil.ReadFile("cmd/server/apigateway.json"); err == nil {
		if err = json.Unmarshal(raw, &apis.APIMap); err == nil {
			a.apis = &apis
			return nil
		}
	}
	log.Println("Loading error: " + err.Error())
	return err
}
