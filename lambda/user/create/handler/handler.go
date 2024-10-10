package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/benjaminkitson/bk-user-api/models"
	utils "github.com/benjaminkitson/bk-user-api/utils/lambda"
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

// TODO: for some error cases, specific messaging would be ideal
// TODO: probably incorporate some sort of request body validation prior to calling cognito or whichever auth provider

func (handler handler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bodyMap := make(map[string]string)

	err := json.Unmarshal([]byte(request.Body), &bodyMap)
	if err != nil {
		handler.logger.Error("Error parsing request body", zap.Error(err))
		return utils.RESPONSE_500, fmt.Errorf("error parsing request body")
	}

	u, err := handler.createUser(ctx, bodyMap)
	if err != nil {
		handler.logger.Error("Failed to get create new user", zap.Error(err))
		return utils.RESPONSE_500, nil
	}
	r, err := json.Marshal(u)
	if err != nil {
		return utils.RESPONSE_500, nil
	}
	return utils.RESPONSE_200(string(r)), nil

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
