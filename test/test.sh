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
      "Description": "GET Request",
      "Request": {
        "RequestVerb":"GET",
        "RequestBody":"",
        "RequestURL":"https://httpbin.org/get"
        }
      }' \
  http://localhost:8080/myservice

  # Create second endpoint
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{
      "Name": "testEP2",
      "Description": "POST Request",
      "Request": {
        "RequestVerb":"POST",
        "RequestBody":"",
        "RequestURL":"https://httpbin.org/post"
        }
      }' \
  http://localhost:8080/myservice

# Get Endpoint
curl \
  --request GET \
  http://localhost:8080/myservice/testEP

# Execute endpoint
curl \
  --request CONNECT \
  http://localhost:8080/myservice/testEP

# Execute all endpoints in API
curl \
  --request CONNECT \
  http://localhost:8080/myservice

