# go-api-gateway
This tool enables interaction between services by creating REST endpoints associated with pre-defined HTTP requests.

### Getting Started
1. Start server on http://localhost:8080: `go run cmd/server/main.go`
2. Create a REST API `myservice`: 
```
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"myservice"}' \
  http://localhost:8080 
```
3. Create an API Endpoint `myservice/testEP` that makes a `GET` request to httpbin.org/get: 
```
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{
      "Name": "testEP",
      "Description": "My EP",
      "Request": {
        "RequestVerb":"GET",
        "RequestBody":"",
        "RequestURL":"https://httpbin.org/get"
        }
      }' \
  http://localhost:8080/myservice
  ```
  4. Execute the endpoint: `curl --request POST http://localhost:8080/myservice/testEP`
