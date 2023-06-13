package service

import (
	"github.com/rs/zerolog/log"

	"github.com/dock-tech/lambda-iam-autentication/internal/core"
	"github.com/dock-tech/lambda-iam-autentication/internal/adapter/restapi"

)

var childLogger = log.With().Str("SERVICE", "AutenticationService").Logger()

type AutenticationService struct {
	adapterRestApi *restapi.AdapterRestApi
}

func NewAutenticationService(adapterRestApi *restapi.AdapterRestApi) *AutenticationService{
	childLogger.Debug().Msg("NewAutenticationService")

	return &AutenticationService{
		adapterRestApi: adapterRestApi,
	}
}

func (s *AutenticationService) AutenticationIAM(autentication core.Autentication) (*core.Autentication, error){
	childLogger.Debug().Msg("AutenticationIAM")

	_, err := s.adapterRestApi.PostAutentication(&autentication)
	if err != nil {
		return nil, err
	}

	return nil, nil
}