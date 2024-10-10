package handler

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/benjaminkitson/bk-user-api/models"
	"go.uber.org/zap"
)

type mockUserStore struct {
	isError bool
}

func (m mockUserStore) Delete(ctx context.Context, id string) (user models.User, err error) {
	return models.User{}, nil
}

/*
Tests the basic workings of the handler
*/
func TestHandler(t *testing.T) {
	type test struct {
		Name                   string
		StoreError             bool
		RequestBody            string
		RequestPath            string
		ExpectedStatusCode     int
		IsHandlerErrorExpected bool
	}

	tests := []test{
		{
			Name:               "Successfully delete user",
			RequestBody:        "{\"email\": \"abc@gmail.com\"}",
			ExpectedStatusCode: 200,
		},
		{
			Name:                   "Failed to create user",
			RequestBody:            "{\"email\": \"abc@gmail.com\"}",
			ExpectedStatusCode:     500,
			StoreError:             true,
			IsHandlerErrorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			l, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("Failed to initialise dev logger")
			}

			u := mockUserStore{
				isError: tt.StoreError,
			}

			h, err := NewHandler(l, u)
			if err != nil {
				t.Fatalf("Failed to initialise handler")
			}

			req := events.APIGatewayProxyRequest{
				// This test should probably fail if the body isn't the correct format?
				Body: tt.RequestBody,
				Path: tt.RequestPath,
			}

			r, err := h.Handle(context.Background(), req)
			if err != nil && !tt.IsHandlerErrorExpected {
				t.Fatalf("Unexpected handler error")
			}

			if r.StatusCode != tt.ExpectedStatusCode {
				t.Fatalf("Expected Status Code to be %v", tt.ExpectedStatusCode)
			}
		})
	}
}
