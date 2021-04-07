package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/corinz/apigateway/internal/app"
	agw "github.com/corinz/apigateway/pkg/apigateway"
)

// TODO, use json as strings, eliminate need to unmarshal objects and test struct fields
// TODO go testing only works with -run flag, go test . and go test api_test.go return different errors
var URL string = "http://localhost:8080"

func TestWebServerRunning(t *testing.T) {
	app := app.NewApp()
	go app.Startup(":8080")
	defer app.Shutdown()
	_, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Error(err)
	}
}

func TestAPICreate(t *testing.T) {
	app := app.NewApp()
	go app.Startup(":8080")
	defer app.Shutdown()

	var json = []byte(`{"Name":"testAPI"}`)
	err := makeJSONReq(json, "/", "testAPI", true)
	if err != nil {
		t.Error(err)
	}
}

func TestAPIEndpoint(t *testing.T) {
	app := app.NewApp()
	go app.Startup(":8080")
	defer app.Shutdown()

	// API Create
	var json = []byte(`{"Name":"testAPI"}`)
	err := makeJSONReq(json, "/", "testAPI", true)
	if err != nil {
		t.Error(err)
	}

	// API Endpoint create
	json = []byte(`{"Name":"testEP","Description":"GET Request","Request":{"RequestVerb":"GET","RequestBody":"","RequestURL":"https://httpbin.org/get"}}`)
	err = makeJSONReq(json, "/testAPI", "testEP", false)
	if err != nil {
		t.Error(err)
	}
}

func TestBadCreate(t *testing.T) {
	app := app.NewApp()
	go app.Startup(":8080")
	defer app.Shutdown()

	// Existing API
	var json = []byte(`{"Name":"testAPI"}`)
	err := makeJSONReq(json, "/", "testAPI", true)
	err = makeJSONReq(json, "/", "testAPI", true)
	if err == nil {
		t.Error(err)
	}

	json = []byte(`{"Name": "testEP","Description": "My EP","HTTPVerb": "POST","JSONPayload": "sleep 5"}`)
	err = makeJSONReq(json, "/testAPI", "testEP", false)
	err = makeJSONReq(json, "/testAPI", "testEP", false)
	if err == nil {
		t.Error(err)
	}
}

// getReq gets an endpoint and tests for a valid decoded json object
func getReq(url string) error {
	resp, err := http.Get(URL + url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	} else if json.Valid(body) {
		return errors.New("json invalid")
	}

	var a agw.API // any struct with 'Name' field can be used
	json.Unmarshal(body, &a)

	if a.Name == "" {
		return errors.New("parm does not exist")
	}

	return nil

}

// makeJSONReq issues a request to the api at a given endpoint
//  and tests the response to ensure a valid struct has been created
func makeJSONReq(jsonStr []byte, endpoint string, nameTest string, mode bool) error {

	// Build request
	req, err := http.NewRequest("POST", URL+endpoint, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	// Do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check http status
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	// Unmarshal response body and test struct in http response
	if mode == true { // api struct
		a, err := readJSONAPI(resp)
		if err != nil {
			return err
		}
		if a.Name != nameTest {
			return errors.New("server returned invalid struct" + a.Name + " != " + nameTest)
		}
	} else { // apiEndpoint struct
		a, err := readJSONAPIEndpoint(resp)
		if err != nil {
			return err
		}
		if a.Name != nameTest {
			return errors.New("server returned invalid struct " + a.Name + " != " + nameTest)
		}
	}

	return nil
}

// readJSONAPI unmarshals byte slice to api struct
func readJSONAPI(resp *http.Response) (agw.API, error) {
	var a agw.API

	// Get body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return a, err
	}
	json.Unmarshal(body, &a)
	return a, nil
}

// readJSONAPIEndpoint unmarshals byte slice to apiEndpoint struct
func readJSONAPIEndpoint(resp *http.Response) (agw.APIEndpoint, error) {
	var a agw.APIEndpoint

	// Get body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return a, err
	}
	json.Unmarshal(body, &a)
	return a, nil
}

// TODO go routine with a write lock and another go. create a wait group to ensure both go routines complete before finishing test
