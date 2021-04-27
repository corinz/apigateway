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

type createAPIResponse struct {
	Name string `json:"Name"`
}

// swagger:route POST /api create createAPI
// Use this request to create an endpoint.
// responses:
//   200: createAPIResponse

// Response Body
// swagger:response createAPIResponse
type createAPIResponseWrapper struct {
	// in:body
	Body agw.API
}

// swagger:parameters createAPI
type foobarParamsWrapper struct {
	// Start by giving your endpoint a na.
	// in:body
	Body agw.API
}
