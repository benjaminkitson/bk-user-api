package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StackProps struct {
	awscdk.StackProps
}

func NewDefaultLambdaProps(path string) *awslambdago.GoFunctionProps {
	return &awslambdago.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		Description:  jsii.String("Handler for user API"),
		Tracing:      awslambda.Tracing_ACTIVE,
		Bundling: &awslambdago.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-trimpath -buildvcs=false`),
		},
		Runtime:    awslambda.Runtime_PROVIDED_AL2(),
		Entry:      jsii.String(path),
		MemorySize: jsii.Number(256),
		Timeout:    awscdk.Duration_Minutes(jsii.Number(1)),
	}
}

func NewStack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	fallbackLambdaProps := NewDefaultLambdaProps("../lambda/fallback")
	fallbackLambda := awslambdago.NewGoFunction(stack, jsii.String("fallbackHandler"), fallbackLambdaProps)

	createUserLambdaProps := NewDefaultLambdaProps("../lambda/user/create")
	createUserLambda := awslambdago.NewGoFunction(stack, jsii.String("createUserHandler"), createUserLambdaProps)

	deleteUserLambdaProps := NewDefaultLambdaProps("../lambda/user/delete")
	deleteUserLambda := awslambdago.NewGoFunction(stack, jsii.String("deleteUserHandler"), deleteUserLambdaProps)

	userDB := awsdynamodb.NewTable(stack, jsii.String("userTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("_pk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("_sk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName:   jsii.String("userTable"),
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	userDB.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("email"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	userDB.GrantReadWriteData(createUserLambda)
	userDB.GrantReadWriteData(deleteUserLambda)

	userApi := awsapigateway.NewLambdaRestApi(stack, jsii.String("Endpoint"), &awsapigateway.LambdaRestApiProps{
		DomainName: &awsapigateway.DomainNameOptions{
			DomainName: jsii.String("api.benjaminkitson.com"),
			Certificate: awscertificatemanager.Certificate_FromCertificateArn(
				stack,
				jsii.String("benjaminkitson-certificate"),
				jsii.String("arn:aws:acm:eu-west-2:905418429454:certificate/42197bf4-d86d-404a-87a6-748c4858d916"),
			),
		},
		DisableExecuteApiEndpoint: jsii.Bool(true),
		RestApiName:               jsii.String("bk-api"),
		Handler:                   fallbackLambda,
		Proxy:                     jsii.Bool(false),
	})

	users := userApi.Root().AddResource(jsii.String("user"), &awsapigateway.ResourceOptions{})

	createUser := users.AddResource(jsii.String("create"), &awsapigateway.ResourceOptions{})
	createUser.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(createUserLambda, &awsapigateway.LambdaIntegrationOptions{}), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_IAM,
	})

	deleteUser := users.AddResource(jsii.String("delete"), &awsapigateway.ResourceOptions{})
	deleteUser.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(deleteUserLambda, &awsapigateway.LambdaIntegrationOptions{}), &awsapigateway.MethodOptions{
		AuthorizationType: awsapigateway.AuthorizationType_IAM,
	})

	z := awsroute53.HostedZone_FromLookup(stack, jsii.String("zone"), &awsroute53.HostedZoneProviderProps{
		DomainName: jsii.String("benjaminkitson.com"),
	})

	awsroute53.NewARecord(stack, jsii.String("apiRecord"), &awsroute53.ARecordProps{
		Zone:       z,
		RecordName: jsii.String("api"),
		Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewApiGateway(userApi)),
	})

	// TODO: Eventually ascertain if the below is really needed
	// awscloudfront.NewDistribution(stack, jsii.String("myDist"), &awscloudfront.DistributionProps{
	// 	DefaultBehavior: &awscloudfront.BehaviorOptions{
	// 		Origin: awscloudfrontorigins.NewRestApiOrigin(pokedexApi, &awscloudfrontorigins.RestApiOriginProps{}),
	// 	},
	// })

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewStack(app, "ApiTestStack", &StackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
