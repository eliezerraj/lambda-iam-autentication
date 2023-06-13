package service

import(
	"testing"
	"github.com/rs/zerolog"

	"github.com/dock-tech/lambda-iam-autentication/internal/adapter/restapi"
	"github.com/dock-tech/lambda-iam-autentication/internal/core"
)

var (
	logLevel	=	zerolog.DebugLevel // InfoLevel DebugLevel
	version		=	"v 0.1"
	squad		=	"architecture"
	AWS_REGION  = 	"us-east-2"

	autenticationAdapterRestApi	*restapi.AdapterRestApi
	restApiData					core.RestApiData
	autentication				core.Autentication
)

func TestAutenticationAIM(t *testing.T) {
	t.Setenv("AWS_REGION", AWS_REGION)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	restApiData.Host = "https://kfyn94nf42.execute-api.us-east-2.amazonaws.com"
	restApiData.Path =  "/live/person"
	autentication.ClientID = "Eliezer"
	autentication.CLientSecret = "1234"
	autentication.ApiKeyID = "xzy"

	autenticationAdapterRestApi = restapi.NewAdapterRestApi(&restApiData)
	autenticationService 		:= NewAutenticationService(autenticationAdapterRestApi)
	
	autenticationService.AutenticationIAM(autentication)

}
