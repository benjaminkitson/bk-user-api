package testhelpers

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBTester struct {
	testClient *dynamodb.Client
}

type testUser struct {
	Username string
}

func (d DBTester) CreateLocalTable(t *testing.T, tableName string) string {
	_, err := d.GetTestClient().CreateTable(context.Background(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("_pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("_sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("_pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("_sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		TableName:   aws.String(tableName),
	})
	if err != nil {
		t.Fatalf("failed to create local table: %v", err)
	}

	userID := "12345"

	testUser := testUser{
		Username: "benk13",
	}

	item, err := attributevalue.MarshalMap(testUser)
	if err != nil {
		t.Fatalf("an error ocurred marshaling the record: %v", err)
	}

	item["_pk"] = &types.AttributeValueMemberS{Value: d.getUserPK(userID)}
	item["_sk"] = &types.AttributeValueMemberS{Value: d.getUserSK(userID)}

	_, err = d.GetTestClient().PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: &tableName,
	})
	if err != nil {
		t.Fatalf("an error ocurred creating the user record: %v", err)
	}

	return tableName
}

func (store DBTester) getUserPK(policyId string) (_pk string) {
	return fmt.Sprintf("user/%s", policyId)
}

func (store DBTester) getUserSK(policyId string) (_pk string) {
	return policyId
}

func (d DBTester) DeleteLocalTable(t *testing.T, name string) {
	_, err := d.GetTestClient().DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(name),
	})
	if err != nil {
		t.Fatalf("failed to delete table: %v", err)
	}
}

func (d DBTester) GetTestClient() *dynamodb.Client {
	if d.testClient != nil {
		return d.testClient
	}
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}
	creds := credentials.NewStaticCredentialsProvider("fake", "accessKeyId", "secretKeyId")
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("eu-pluto-1"),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load test aws config %w", err))
	}
	d.testClient = dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolverV2(dynamodb.NewDefaultEndpointResolverV2()), func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	})
	return d.testClient
}
