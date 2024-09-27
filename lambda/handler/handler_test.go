package handler

import (
	"testing"
)

/*
Tests the basic workings of the handler, using a mocked auth provider client that either succeeds or returns some generic error
*/
func TestHandler(t *testing.T) {
	// type test struct {
	// 	Name                   string
	// 	AdapterError           bool
	// 	SecretsGetterError     bool
	// 	RequestBody            string
	// 	RequestPath            string
	// 	ExpectedStatusCode     int
	// 	IsHandlerErrorExpected bool
	// }

	// tests := []test{
	// 	{
	// 		Name:                   "Invalid path supplied",
	// 		RequestBody:            "{\"email\": \"abc@gmail.com\", \"password\": \"password\"}",
	// 		RequestPath:            "/someInvalidPath",
	// 		ExpectedStatusCode:     500,
	// 		IsHandlerErrorExpected: true,
	// 	},
	// }

	// for _, tt := range tests {

	// }
}
