#!/bin/sh

# Create myservice API
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"myservice"}' \
  http://localhost:8080 

# Get list of all APIs
curl \
  --request GET \
  http://localhost:8080 

# Create endpoint
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

# Get Endpoint
curl \
  --request GET \
  http://localhost:8080/myservice/testEP

# Execute endpoint
curl \
  --request POST \
  http://localhost:8080/myservice/testEP

