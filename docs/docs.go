// Package classification API Gateway.
//
// Documentation of API Gateway.
//
//     Schemes: https
//     BasePath: /
//     Version: 1.0.0
//     Host: localhost
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package docs

import agw "github.com/corinz/apigateway/pkg/apigateway"

//
// Params
//

//
// swagger:parameters createAPIEndpoint
// swagger:parameters listAPI
// swagger:parameters listAPIEndpoint
// swagger:parameters executeAPI
// swagger:parameters executeAPIEndpoint
// swagger:parameters deleteAPI
// swagger:parameters deleteAPIEndpoint
type APIName struct {
	// type: string
	// In: path
	APIName string
}

//
// swagger:parameters listAPIEndpoint
// swagger:parameters executeAPIEndpoint
// swagger:parameters deleteAPIEndpoint
type APIEndpoint struct {
	// type: string
	// In: path
	APIEndpoint string
}

//
//POSTs
//

//
// swagger:route POST /api create createAPI
// Use this request to create an API.
// responses:
//   200: createAPIResponse

// Response Body
// swagger:response createAPIResponse
type createAPIResponseWrapper struct {
	// in:body
	Body agw.API
}

type createAPIResponse struct {
	Name string `json:"Name"`
}

// swagger:parameters createAPI
type createAPIParamsWrapper struct {
	// Give your API a name.
	// in:body
	Body createAPIResponse
}

//
// swagger:route POST /api/{APIName} create createAPIEndpoint
// Use this request to create an API Endpoint using the /api generated API Name.
// responses:
//   200: createAPIEndpointResponse

// Response Body
// swagger:response createAPIEndpointResponse
type createAPIEndpointResponseWrapper struct {
	// in:body
	Body agw.APIEndpoint
}

type createAPIEndpointResponse struct {
	Name        string
	Description string
	Request     struct {
		RequestBody string
		RequestURL  string
		RequestVerb string
	}
}

// swagger:parameters createAPIEndpoint
type createAPIEndpointParamsWrapper struct {
	// Give your API Endpoint a name.
	// in:body
	Body createAPIEndpointResponse
}

//
// GETs
//

//
// swagger:route GET /api list listAPIs
// Use this request to list all APIs.
// responses:
//   200: listAPIs

// Response Body
// swagger:response listAPIs
type listAPIs struct {
	// in:body
	Body agw.APIs
}

//
// swagger:route GET /api/{APIName} list listAPI
// Use this request to list an API.
// responses:
//   200: listAPIResponse

// Response Body
// swagger:response listAPIResponse
type listAPIResponseWrapper struct {
	// in:body
	Body agw.API
}

//
// swagger:route GET /api/{APIName}/{APIEndpoint} list listAPIEndpoint
// Use this request to list an API Endpoint.
// responses:
//   200: listAPIEndpointResponse

// Response Body
// swagger:response listAPIEndpointResponse
type listAPIEndpointResponseWrapper struct {
	// in:body
	Body agw.APIEndpoint
}

//
// swagger:route GET /{APIName} execute executeAPI
// Use this request to execute all Endpoints in an API.
// responses:
//   200: executeAPIResponse

// Response Body
// swagger:response executeAPIResponse
type executeAPIResponseWrapper struct {
	// in:body
	Body agw.API
}

//
// swagger:route GET /{APIName}/{APIEndpoint} execute executeAPIEndpoint
// Use this request to execute an API Endpoint.
// responses:
//   200: executeAPIEndpointResponse

// Response Body
// swagger:response executeAPIEndpointResponse
type executeAPIEndpointResponseWrapper struct {
	// in:body
	Body agw.APIEndpoint
}

//
// DELETEs
//

//
// swagger:route DELETE /api/{APIName} delete deleteAPI
// Use this request to delete an API and all of its endpoints.
// responses:
//   200: deleteAPIResponse

// Response Body
// swagger:response deleteAPIResponse
type deleteAPIResponseWrapper struct {
	// in:body
	Body agw.API
}

//
// swagger:route DELETE /api/{APIName}/{APIEndpoint} delete deleteAPIEndpoint
// Use this request to delete an Endpoint within an API.
// responses:
//   200: deleteAPIEndpointResponse

// Response Body
// swagger:response deleteAPIEndpointResponse
type deleteAPIEndpointResponseWrapper struct {
	// in:body
	Body agw.APIEndpoint
}
