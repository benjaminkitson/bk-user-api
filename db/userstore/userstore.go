package userstore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/benjaminkitson/bk-user-api/models"
	"github.com/pkg/errors"
)

const (
	PKKey   string = "_pk"
	GSI1Key string = "_gsi1"
	// SKKey string = "_sk"
)

type UserStore struct {
	tableName string
	client    *dynamodb.Client
}

func NewUserStore(client *dynamodb.Client, tableName string) UserStore {
	return UserStore{
		tableName: tableName,
		client:    client,
	}
}

func (store UserStore) GetByID(ctx context.Context, id string) (models.User, error) {
	key := map[string]types.AttributeValue{
		PKKey: &types.AttributeValueMemberS{Value: store.getUserPK(id)},
		// SKKey: &types.AttributeValueMemberS{Value: store.getUserSK(id)},
	}
	query := dynamodb.GetItemInput{
		TableName:      &store.tableName,
		Key:            key,
		ConsistentRead: aws.Bool(true),
	}

	item, err := store.client.GetItem(ctx, &query)
	if err != nil {
		return models.User{}, err
	}

	if len(item.Item) == 0 {
		return models.User{}, err
	}

	var user models.User
	err = attributevalue.UnmarshalMap(item.Item, &user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (store UserStore) GetByEmail(ctx context.Context, email string) (models.User, error) {
	out, err := store.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &store.tableName,
		IndexName:              aws.String("gsi1"),
		KeyConditionExpression: aws.String("#gsi1 = :gsi1"),
		ExpressionAttributeNames: map[string]string{
			"#gsi1": "_gsi1",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi1": &types.AttributeValueMemberS{Value: store.getUserGSI1(email)},
		},
	})
	if err != nil {
		return models.User{}, err
	}

	if len(out.Items) > 1 {
		return models.User{}, fmt.Errorf("expected maximum of 1 records, found %d", len(out.Items))
	}

	if len(out.Items) == 0 {
		return models.User{}, nil
	}

	var user models.User
	err = attributevalue.UnmarshalMap(out.Items[0], &user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (store UserStore) Put(ctx context.Context, record models.User) (models.User, error) {
	item, err := attributevalue.MarshalMap(record)
	if err != nil {
		return models.User{}, errors.Wrap(err, "an error ocurred marshaling the record")
	}

	item[PKKey] = &types.AttributeValueMemberS{Value: store.getUserPK(record.UserID)}
	item[GSI1Key] = &types.AttributeValueMemberS{Value: store.getUserPK(record.Email)}
	// item[SKKey] = &types.AttributeValueMemberS{Value: store.getUserSK(id)}

	_, err = store.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: &store.tableName,
	})

	if err != nil {
		return models.User{}, err
	}

	return record, err
}

func (store UserStore) Delete(ctx context.Context, id string) (string, error) {
	// TODO: ascertain if these attributes are needed
	_, err := store.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(store.tableName),
		Key: map[string]types.AttributeValue{
			PKKey: &types.AttributeValueMemberS{Value: store.getUserPK(id)},
			// SKKey: &types.AttributeValueMemberS{Value: store.getUserSK(id)},
		},
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (store UserStore) getUserPK(userID string) (_pk string) {
	return fmt.Sprintf("user/%s", userID)
}

func (store UserStore) getUserGSI1(email string) (gsi1 string) {
	return fmt.Sprintf("email/%s", email)
}

func (store UserStore) getUserSK(userID string) (_pk string) {
	return userID
}
