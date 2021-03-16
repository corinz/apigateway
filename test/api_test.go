package APIGateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

// TODO, use json as strings, eliminate need to unmarshal objects and test struct fields

var URL string = "http://localhost:8080"

func TestWebServerRunning(t *testing.T) {
	go Startup()
	_, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Error()
	}
}

func TestAPICreate(t *testing.T) {
	go Startup()
	var json = []byte(`{"Name":"testAPI"}`)

	// test for true negative error
	err := makeJSONReq(json, "/", "testAPI", true)
	if err != nil {
		t.Error()
	}
}

func TestAPIEndpoint(t *testing.T) {
	go Startup()

	// API Create
	var json = []byte(`{"Name":"testAPI"}`)
	err := makeJSONReq(json, "/", "testAPI", true)
	if err != nil {
		t.Error()
	}

	// API Endpoint create
	json = []byte(`{"Name": "testEP","Description": "My EP","HTTPVerb": "POST","Command": "sleep 5"}`)

	// test for true positive error
	err = makeJSONReq(json, "/testAPI", "testEP", false)
	if err != nil {
		t.Error()
	}
}

func TestBadCreate(t *testing.T) {
	go Startup()

	// Existing API
	var json = []byte(`{"Name":"testAPI"}`)
	err := makeJSONReq(json, "/", "testAPI", true)
	err = makeJSONReq(json, "/", "testAPI", true)
	if err == nil {
		t.Error()
	}

	json = []byte(`{"Name": "testEP","Description": "My EP","HTTPVerb": "POST","Command": "sleep 5"}`)
	err = makeJSONReq(json, "/testAPI", "testEP", false)
	err = makeJSONReq(json, "/testAPI", "testEP", false)
	if err == nil {
		t.Error()
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

	var a api // any struct with 'Name' field can be used
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
			return errors.New("server returned invalid struct")
		}
	} else { // apiEndpoint struct
		a, err := readJSONAPIEndpoint(resp)
		if err != nil {
			return err
		}
		if a.Name != nameTest {
			return errors.New("server returned invalid struct")
		}
	}

	return nil
}

// readJSONAPI unmarshals byte slice to api struct and tests the 'Name' field
func readJSONAPI(resp *http.Response) (api, error) {
	var a api

	// Get body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return a, err
	}
	json.Unmarshal(body, &a)
	return a, nil
}

// readJSONAPIEndpoint unmarshals byte slice to apiEndpoint struct and tests the 'Name' field
func readJSONAPIEndpoint(resp *http.Response) (apiEndpoint, error) {
	var a apiEndpoint

	// Get body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return a, err
	}
	json.Unmarshal(body, &a)
	return a, nil
}
