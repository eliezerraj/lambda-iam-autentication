package main

import(
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dock-tech/lambda-iam-autentication/internal/service"
	"github.com/dock-tech/lambda-iam-autentication/internal/adapter/restapi"
	"github.com/dock-tech/lambda-iam-autentication/internal/core"
)

var (
	logLevel	=	zerolog.DebugLevel // InfoLevel DebugLevel
	version		=	"v 0.1"
	squad		=	"architecture"

	autenticationService 		*service.AutenticationService
	autenticationAdapterRestApi	*restapi.AdapterRestApi
	restApiData					core.RestApiData
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	// mock
	restApiData.Host = "https://kfyn94nf42.execute-api.us-east-2.amazonaws.com"
	restApiData.Path =  "/live/person"
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

	if os.Getenv("HOST_AUTH") !=  "" {
		restApiData.Host = os.Getenv("HOST_AUTH")
	}
	if os.Getenv("PATH_AUTH") !=  "" {
		restApiData.Path = os.Getenv("PATH_AUTH")
	}
}

func main() {
	log.Debug().Msg("main lambda-iam-autentication")
	log.Debug().Str("version", version).Msg("")

	autenticationAdapterRestApi = restapi.NewAdapterRestApi(&restApiData)
	autenticationService = service.NewAutenticationService(autenticationAdapterRestApi)

	///autenticationService.AutenticationIAM(nil)
}