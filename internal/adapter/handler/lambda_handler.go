package handler

import(
	"github.com/rs/zerolog/log"
	"net/http"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-lambda-go/events"

	"github.com/dock-tech/lambda-iam-autentication/internal/service"
	"github.com/dock-tech/lambda-iam-autentication/internal/erro"
	"github.com/dock-tech/lambda-iam-autentication/internal/core"

)

var childLogger = log.With().Str("HANDLER", "IamAutenticationHandler").Logger()

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

type MessageBody struct {
	Msg *string `json:"message,omitempty"`
}

type IamAutenticationHandler struct {
	autenticationService 	service.AutenticationService
}

func NewIamAutenticationHandler(autenticationService service.AutenticationService) *IamAutenticationHandler{
	childLogger.Debug().Msg("NewIamAutenticationHandler")

	return &IamAutenticationHandler{
		autenticationService: autenticationService,
	}
}

func (h *IamAutenticationHandler) UnhandledMethod() (*events.APIGatewayProxyResponse, error){
	return ApiHandlerResponse(http.StatusMethodNotAllowed, ErrorBody{aws.String(erro.ErrMethodNotAllowed.Error())})
}

func (h *IamAutenticationHandler) AutenticationIAM(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("GetFinancialMovimentByPerson")

    var autentication core.Autentication
    if err := json.Unmarshal([]byte(req.Body), &autentication); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
    }

	response, err := h.autenticationService.AutenticationIAM(autentication)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *IamAutenticationHandler) GetVersion(version string) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("GetVersion")

	response := MessageBody { Msg: &version }
	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	
	return handlerResponse, nil
}
