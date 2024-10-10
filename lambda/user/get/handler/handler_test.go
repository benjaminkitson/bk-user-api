package handler

import (
	"context"
	"fmt"

	"github.com/benjaminkitson/bk-user-api/models"
)

type mockUserStore struct {
	isError bool
}

func (m mockUserStore) GetByID(ctx context.Context, id string) (user models.User, err error) {
	return models.User{}, nil
}

// TODO: when a user already exists, GetByEmail error
func (m mockUserStore) GetByEmail(ctx context.Context, email string) (user models.User, err error) {
	return models.User{}, nil
}

func (m mockUserStore) Put(ctx context.Context, record models.User, id string) (models.User, error) {
	if m.isError {
		return models.User{}, fmt.Errorf("UserStore put error!")
	}
	return record, nil
}

/*
Tests the basic workings of the handler
*/
// func TestHandler(t *testing.T) {
// 	type test struct {
// 		Name                   string
// 		StoreError             bool
// 		RequestBody            string
// 		RequestPath            string
// 		ExpectedStatusCode     int
// 		IsHandlerErrorExpected bool
// 	}

// 	tests := []test{
// 		{
// 			Name:               "Successfully create user",
// 			RequestBody:        "{\"email\": \"abc@gmail.com\"}",
// 			RequestPath:        "/user/create",
// 			ExpectedStatusCode: 200,
// 		},
// 		{
// 			Name:                   "Failed to create user",
// 			RequestBody:            "{\"email\": \"abc@gmail.com\"}",
// 			RequestPath:            "/user/create",
// 			ExpectedStatusCode:     500,
// 			StoreError:             true,
// 			IsHandlerErrorExpected: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			l, err := zap.NewDevelopment()
// 			if err != nil {
// 				t.Fatalf("Failed to initialise dev logger")
// 			}

// 			u := mockUserStore{
// 				isError: tt.StoreError,
// 			}

// 			h, err := NewHandler(l, u)
// 			if err != nil {
// 				t.Fatalf("Failed to initialise handler")
// 			}

// 			req := events.APIGatewayProxyRequest{
// 				// This test should probably fail if the body isn't the correct format?
// 				Body: tt.RequestBody,
// 				Path: tt.RequestPath,
// 			}

// 			r, err := h.Handle(context.Background(), req)
// 			if err != nil && !tt.IsHandlerErrorExpected {
// 				t.Fatalf("Unexpected handler error")
// 			}

// 			if r.StatusCode != tt.ExpectedStatusCode {
// 				t.Fatalf("Expected Status Code to be %v", tt.ExpectedStatusCode)
// 			}
// 		})
// 	}
// }
