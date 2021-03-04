#!/bin/sh

# Create myservice API
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"myservice","route":"myservice"}' \
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
      "HTTPVerb": "POST",
      "Command": "sleep 30"
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

