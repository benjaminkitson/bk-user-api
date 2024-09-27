package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/benjaminkitson/user-api/db/userstore"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type handler struct {
	logger    *zap.Logger
	userStore userstore.UserStore
}

func NewHandler(logger *zap.Logger) (handler, error) {
	return handler{
		logger: logger,
	}, nil
}

var Headers = map[string]string{
	"Access-Control-Allow-Headers": "Content-Type",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
}

const GenericError = "{\"message\": \"Something went wrong!\"}"

// TODO: make distinction between 400 and 500 errors
// TODO: understand how different methods are dealt with (post vs get etc)
// TODO: probably incorporate some sort of request body validation prior to calling cognito or whichever auth provider

func (handler handler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bodyMap := make(map[string]string)

	err := json.Unmarshal([]byte(request.Body), &bodyMap)
	if err != nil {
		handler.logger.Error("Error parsing request body", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    Headers,
			// Error body needed? Probably not
		}, fmt.Errorf("error parsing request body")
	}

	if request.Path == "/create" {
		err := handler.createUser(ctx, bodyMap)
		if err != nil {
			handler.logger.Error("Failed to get create new user", zap.Error(err))
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    Headers,
				Body:       GenericError,
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    Headers,
		}, nil
	}

	handler.logger.Error("invalid path")
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers:    Headers,
		// Error body needed? Probably not
	}, fmt.Errorf("invalid path")
}

func (handler handler) createUser(ctx context.Context, requestBody map[string]string) error {
	id := uuid.New().String()

	user := userstore.User{
		Username: requestBody["username"],
	}

	err := handler.userStore.Put(ctx, user, id)
	return err
}
