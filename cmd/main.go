package main

import(
	"os"
	"fmt"
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/dock-tech/lambda-iam-autentication/internal/service"
	"github.com/dock-tech/lambda-iam-autentication/internal/adapter/restapi"
	"github.com/dock-tech/lambda-iam-autentication/internal/core"
	"github.com/dock-tech/lambda-iam-autentication/internal/adapter/handler"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

)

var (
	logLevel	=	zerolog.DebugLevel // InfoLevel DebugLevel
	version		string
	squad		string
	awsSecretId string
	env			string

	xApigwApiId string
	autenticationService 		*service.AutenticationService
	autenticationAdapterRestApi	*restapi.AdapterRestApi
	restApiData					core.RestApiData
	autenticationData			core.Autentication
	response					*events.APIGatewayProxyResponse
	autenticationHandler		*handler.IamAutenticationHandler
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	// mock data
	restApiData.Host = "https://kfyn94nf42.execute-api.us-east-2.amazonaws.com"
	restApiData.Path =  "/live/person"
	awsSecretId = "secretInternalApp"
	version = "v 0.1"
	squad = "ARCH"
	env = "DEV"
	xApigwApiId = "qweqwewqew" 
	//
	getEnv()
}

func getEnv() {
	if os.Getenv("LOG_LEVEL") !=  "" {
		if (os.Getenv("LOG_LEVEL") == "DEBUG"){
			logLevel = zerolog.DebugLevel
		}else if (os.Getenv("LOG_LEVEL") == "INFO"){
			logLevel = zerolog.InfoLevel
		}else if (os.Getenv("LOG_LEVEL") == "ERROR"){
				logLevel = zerolog.ErrorLevel
		}else {
			logLevel = zerolog.DebugLevel
		}
	}
	if os.Getenv("VERSION") !=  "" {
		version = os.Getenv("VERSION")
	}
	if os.Getenv("SQUAD") !=  "" {
		squad = os.Getenv("SQUAD")
	}
	if os.Getenv("ENV") !=  "" {
		squad = os.Getenv("ENV")
	}

	if os.Getenv("HOST_AUTH") !=  "" {
		restApiData.Host = os.Getenv("HOST_AUTH")
	}
	if os.Getenv("PATH_AUTH") !=  "" {
		restApiData.Path = os.Getenv("PATH_AUTH")
	}
	if os.Getenv("AWS_SECRET_ID") !=  "" {
		awsSecretId = os.Getenv("AWS_SECRET_ID")
	}
	if os.Getenv("X_APIGW_API_ID") !=  "" {
		xApigwApiId = os.Getenv("X_APIGW_API_ID")
	}
}

func main() {
	log.Debug().Msg("main lambda-iam-autentication")
	log.Debug().Str("version", version).Msg("")

	// Load the IAM Secret
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("error LoadDefaultConfig")
		return
	}

	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(awsSecretId),
		VersionStage: aws.String("AWSCURRENT"),
	}
	
	secret_result, err := svc.GetSecretValue(context.TODO() ,input)
	if err != nil {
		log.Error().Err(err).Msg("error GetSecretValue")
		return
	}

	//Set the Autentication data
	json.Unmarshal([]byte(*secret_result.SecretString) , &autenticationData)
	autenticationData.ApiKeyID 		= xApigwApiId

	// Create RESTAPI Adapter
	autenticationAdapterRestApi = restapi.NewAdapterRestApi(&restApiData)
	autenticationService 		= service.NewAutenticationService(autenticationAdapterRestApi, &autenticationData)
	autenticationHandler 		= handler.NewIamAutenticationHandler(*autenticationService)

	fmt.Println(autenticationHandler)
}

func lambdaHandler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debug().Msg("lambdaHandler")

	switch req.HTTPMethod {
		case "GET":
			 if (req.Resource == "/version"){
				response, _ = autenticationHandler.GetVersion(version)
			}else {
				response, _ = autenticationHandler.UnhandledMethod()
			}
		case "POST":
			if (req.Resource == "/financialmovimentbyperson"){
				response, _ = autenticationHandler.AutenticationIAM(req)
			}else {
				response, _ = autenticationHandler.UnhandledMethod()
			}
		case "DELETE":
			response, _ = autenticationHandler.UnhandledMethod()
		case "PUT":
			response, _ = autenticationHandler.UnhandledMethod()
		default:
			response, _ = autenticationHandler.UnhandledMethod()
	}
	return response, nil
}