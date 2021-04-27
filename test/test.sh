#!/bin/sh

# Create myservice API
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"myservice"}' \
  https://localhost/api --insecure

# Get list of all APIs
curl \
  --request GET \
  https://localhost/api --insecure

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
  https://localhost/api/myservice --insecure

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
  https://localhost/api/myservice --insecure

# You can request an API and API Endpoints be created in a single request
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{
    "Name": "myservice2",
    "APIEPs":[{
      "Description": "string",
      "Name": "testEP",
      "ParentName": "string",
      "Request": {
        "RequestBody": "string",
        "RequestURL": "string",
        "RequestVerb": "string"}}]}' \
  https://localhost/api --insecure

# Get Endpoint
curl \
  --request GET \
  https://localhost/api/myservice/testEP --insecure

# Execute endpoint
curl \
  --request GET \
  https://localhost/myservice/testEP --insecure

# Execute all endpoints in API
curl \
  --request GET \
  https://localhost/myservice --insecure