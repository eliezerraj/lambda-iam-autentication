package service

import (
	"github.com/rs/zerolog/log"

	"github.com/dock-tech/lambda-iam-autentication/internal/core"
	"github.com/dock-tech/lambda-iam-autentication/internal/adapter/restapi"

)

var childLogger = log.With().Str("SERVICE", "AutenticationService").Logger()

type AutenticationService struct {
	adapterRestApi 		*restapi.AdapterRestApi
	autenticationData	*core.Autentication
}

func NewAutenticationService(	adapterRestApi 		*restapi.AdapterRestApi,
								autenticationData 	*core.Autentication) *AutenticationService{
	childLogger.Debug().Msg("NewAutenticationService")

	return &AutenticationService{
		adapterRestApi: adapterRestApi,
		autenticationData: autenticationData,
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