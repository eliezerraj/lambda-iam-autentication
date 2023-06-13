package restapi

import(
	"context"
	"fmt"
	"net/http"
	"time"
	"crypto/sha256"
	"encoding/json"

	"github.com/mitchellh/mapstructure"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dock-tech/lambda-iam-autentication/internal/core"
	"github.com/dock-tech/lambda-iam-autentication/internal/erro"
)

var childLogger = log.With().Str("RESTAPI", "RestApi").Logger()

type AdapterRestApi struct {
	restApiAuth *core.RestApiData
}

func NewAdapterRestApi(restApiAuth *core.RestApiData) (*AdapterRestApi){
	childLogger.Debug().Msg("NewAdapterRestApi")

	return &AdapterRestApi {
		restApiAuth: restApiAuth,
	}
}

func (r *AdapterRestApi) PostAutentication(autentication *core.Autentication) (*core.Autentication, error) {
	childLogger.Debug().Msg("PostAutentication")

	url := r.restApiAuth.Host + r.restApiAuth.Path 

	autentication_interface, err := makePostAuthIAM(url, autentication)
	if err != nil {
		childLogger.Error().Err(err).Msg("error makePostAuthIAM")
		return nil, err
	}

	var autentication_result core.Autentication
	err = mapstructure.Decode(autentication_interface, &autentication_result)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, erro.ErrParceInterface
    }
    
	return &autentication_result, nil
}

func makePostAuthIAM(url string, inter interface{}) (interface{}, error) {
	childLogger.Debug().Msg("makePostAuthIAM")
	childLogger.Debug().Str("",url).Msg("")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		childLogger.Error().Err(err).Msg("error get Context")
		return false, erro.ErrBadRequest
	}
	credentials, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		childLogger.Error().Err(err).Msg("error get Credentials")
		return false, erro.ErrBadRequest
	}

	fmt.Println("***** ",credentials," === ",cfg.Region)

	// Set timeout APIGW
	client := &http.Client{Timeout: time.Second * 29}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		childLogger.Error().Err(err).Msg("error NewRequest")
		return false, erro.ErrBadRequest
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8");
	//req.Header.Add("x-api-key", "NvbtUByZfK3PO0trigViL2OSQPlx8KTMa8pszPgt");
	hash := sha256.Sum256([]byte(""))
	hexHash := fmt.Sprintf("%x", hash)

	signer := v4.NewSigner()
	err = signer.SignHTTP(context.TODO(), credentials, req, hexHash, "execute-api", cfg.Region, time.Now())
	if err != nil {
		childLogger.Error().Err(err).Msg("error signer with credentials")
		return false, erro.ErrBadRequest
	}

	resp, err := client.Do(req)
	if err != nil {
		childLogger.Error().Err(err).Msg("error client.Do")
		return false, erro.ErrBadRequest
	}

	childLogger.Debug().Int("StatusCode :", resp.StatusCode).Msg("----")
	switch (resp.StatusCode) {
		case 401:
			return false, erro.ErrHTTPForbiden
		case 403:
			return false, erro.ErrHTTPForbiden
		case 200:
		case 400:
			return false, erro.ErrNotFound
		case 404:
			return false, erro.ErrNotFound
		default:
			return false, erro.ErrHTTPForbiden
	}

	result := inter
	err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
		childLogger.Error().Err(err).Msg("error no ErrUnmarshal")
		return false, erro.ErrUnmarshal
    }

	return result, nil
}