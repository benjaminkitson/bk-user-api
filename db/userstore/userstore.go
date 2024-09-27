package userstore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

const (
	pkKey string = "_pk"
	skKey string = "_sk"
)

type User struct {
	Username string
}

type UserStore struct {
	tableName string
	client    *dynamodb.Client
}

func NewUserStore(client *dynamodb.Client) UserStore {
	return UserStore{
		tableName: "users",
		client:    client,
	}
}

func (store UserStore) Get(ctx context.Context, id string) (user User, err error) {
	key := map[string]types.AttributeValue{
		pkKey: &types.AttributeValueMemberS{Value: store.getUserPK(id)},
		skKey: &types.AttributeValueMemberS{Value: store.getUserSK(id)},
	}
	query := dynamodb.GetItemInput{
		TableName:      &store.tableName,
		Key:            key,
		ConsistentRead: aws.Bool(true),
	}

	item, err := store.client.GetItem(ctx, &query)
	if err != nil {
		return
	}

	if len(item.Item) == 0 {
		user = User{}
		return
	}

	err = attributevalue.UnmarshalMap(item.Item, &user)
	if err != nil {
		return
	}

	return
}

func (store UserStore) Put(ctx context.Context, record User, id string) (err error) {
	item, err := attributevalue.MarshalMap(record)
	if err != nil {
		return errors.Wrap(err, "an error ocurred marshaling the record")
	}

	item[pkKey] = &types.AttributeValueMemberS{Value: store.getUserPK(id)}
	item[skKey] = &types.AttributeValueMemberS{Value: store.getUserSK(id)}

	_, err = store.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: &store.tableName,
	})

	if err != nil {
		return err
	}

	return
}

func (store UserStore) getUserPK(policyId string) (_pk string) {
	return fmt.Sprintf("user/%s", policyId)
}

func (store UserStore) getUserSK(policyId string) (_pk string) {
	return policyId
}
