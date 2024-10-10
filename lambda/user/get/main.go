package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/benjaminkitson/bk-user-api/db/userstore"
	"github.com/benjaminkitson/bk-user-api/lambda/user/get/handler"
	"go.uber.org/zap"
)

func main() {
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		logger, err := zap.NewProduction()
		if err != nil {
			fmt.Printf("Failed to initialise logger: %v", err)
			logger = &zap.Logger{}
		}
		defer logger.Sync()

		sdkConfig, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			logger.Error("Failed to intialise SDK config", zap.Error(err))
			return events.APIGatewayProxyResponse{}, err
		}

		d := dynamodb.NewFromConfig(sdkConfig)

		tableName := "userTable"

		u := userstore.NewUserStore(d, tableName)

		h, err := handler.NewHandler(logger, u)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return h.Handle(ctx, request)
	})
}
