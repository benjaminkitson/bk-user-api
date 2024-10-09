package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/benjaminkitson/bk-user-api/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type handler struct {
	logger    *zap.Logger
	userStore handlerUserStore
}

type handlerUserStore interface {
	GetByID(ctx context.Context, id string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Put(ctx context.Context, record models.User, id string) (models.User, error)
}

func NewHandler(logger *zap.Logger, u handlerUserStore) (handler, error) {
	return handler{
		logger:    logger,
		userStore: u,
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

	if request.Path == "/user/create" {
		u, err := handler.createUser(ctx, bodyMap)
		if err != nil {
			handler.logger.Error("Failed to get create new user", zap.Error(err))
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    Headers,
				Body:       GenericError,
			}, nil
		}
		r, err := json.Marshal(u)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    Headers,
				Body:       GenericError,
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    Headers,
			Body:       string(r),
		}, nil
	}

	handler.logger.Error("invalid path", zap.String("path", request.Path))
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers:    Headers,
		Body:       "{\"message\": \"Invalid path\"}",
	}, nil
}

func (handler handler) createUser(ctx context.Context, requestBody map[string]string) (models.User, error) {
	id := uuid.New().String()

	user := models.User{
		Email: requestBody["email"],
	}

	c, err := handler.userStore.GetByEmail(ctx, requestBody["email"])
	if err != nil {
		handler.logger.Error("error checking for existing user by email")
		return models.User{}, err
	}

	if c.Email != "" {
		handler.logger.Error("user with email already exists")
		return models.User{}, fmt.Errorf("user with email already exists")
	}

	u, err := handler.userStore.Put(ctx, user, id)
	if err != nil {
		return models.User{}, err
	}

	return u, err
}
