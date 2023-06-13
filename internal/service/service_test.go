package service

import(
	"testing"

	"context"
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/dock-tech/lambda-iam-autentication/internal/adapter/restapi"
	"github.com/dock-tech/lambda-iam-autentication/internal/core"
)

var (
	logLevel	=	zerolog.DebugLevel // InfoLevel DebugLevel
	version		=	"v 0.1"
	squad		=	"architecture"
	AWS_REGION  = 	"us-east-2"
	xApigwApiId 	string
	awsSecretId		string

	autenticationAdapterRestApi	*restapi.AdapterRestApi
	restApiData					core.RestApiData
	autenticationData			core.Autentication
)

func TestAutenticationAIM(t *testing.T) {
	t.Setenv("AWS_REGION", AWS_REGION)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	restApiData.Host = "https://kfyn94nf42.execute-api.us-east-2.amazonaws.com"
	restApiData.Path =  "/live/person/007"
	awsSecretId 	= "secretInternalApp"

	// Load the IAM Secret
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Errorf("Error -TestAutenticationAIM-LoadDefaultConfig %v ", err)
	}
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(awsSecretId),
		VersionStage: aws.String("AWSCURRENT"),
	}
	
	secret_result, err := svc.GetSecretValue(context.TODO() ,input)
	if err != nil {
		t.Errorf("Error -TestAutenticationAIM-GetSecretValue %v ", err)
	}
	
	//Set the Autentication data
	json.Unmarshal([]byte(*secret_result.SecretString) , &autenticationData)
	autenticationData.ApiKeyID 		= xApigwApiId

	// Create RESTAPI Adapter
	autenticationAdapterRestApi = restapi.NewAdapterRestApi(&restApiData)
	autenticationService 		:= NewAutenticationService(autenticationAdapterRestApi, &autenticationData)
	
	autenticationService.AutenticationIAM(autenticationData)

}
