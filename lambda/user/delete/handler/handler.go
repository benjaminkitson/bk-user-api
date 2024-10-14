package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	utils "github.com/benjaminkitson/bk-user-api/utils/lambda"
	"go.uber.org/zap"
)

type handler struct {
	logger    *zap.Logger
	userStore handlerUserStore
}

type handlerUserStore interface {
	Delete(ctx context.Context, id string) (string, error)
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

	handler.logger.Info("attempting user deletion", zap.String("userID", bodyMap["id"]))
	id, err := handler.userStore.Delete(ctx, bodyMap["id"])
	if err != nil {
		handler.logger.Error("error deleting user", zap.String("userID", bodyMap["id"]), zap.Error(err))
		return utils.RESPONSE_500, nil
	}
	handler.logger.Info("successfully deleted user from db", zap.String("userID", bodyMap["id"]))

	s := map[string]string{
		"id": id,
	}
	r, err := json.Marshal(s)
	if err != nil {
		return utils.RESPONSE_500, nil
	}

	return utils.RESPONSE_200(string(r)), nil
}
