package restapi

import(
	"context"
	"fmt"
	"net/http"
	"time"
	"crypto/sha256"
	"encoding/json"
	"bytes"

	"github.com/mitchellh/mapstructure"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/dock-tech/lambda-iam-autentication/internal/core"
	"github.com/dock-tech/lambda-iam-autentication/internal/erro"
)

var childLogger = log.With().Str("RESTAPI", "=>").Logger()

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

	// Make a REST CALL 
	autentication_interface, err := makePostAuthIAM(url, autentication)
	if err != nil {
		childLogger.Error().Err(err).Msg("error makePostAuthIAM")
		return nil, err
	}

	//Decode the result
	var autentication_result core.Autentication
	err = mapstructure.Decode(autentication_interface, &autentication_result)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, erro.ErrParceInterface
    }

	// Clean the Keys
	autentication_result.ClientID = ""
	autentication_result.CLientSecret = ""

	return &autentication_result, nil
}

func makePostAuthIAM(url string, inter interface{}) (interface{}, error) {
	childLogger.Debug().Msg("makePostAuthIAM")
	childLogger.Debug().Str("",url).Msg("")
	childLogger.Debug().Interface("==>",inter).Msg("<==")

	ctx := context.TODO()
	autentication_data, ok:= inter.(*core.Autentication)
	if !ok {
		return false, erro.ErrBadRequest
	}
	// Get Credentials informed via ENV
	staticProvider := credentials.NewStaticCredentialsProvider(	autentication_data.ClientID, 
																autentication_data.CLientSecret, 
																"")
	// ---------------------------
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(staticProvider))
	if err != nil {
		childLogger.Error().Err(err).Msg("error get Context")
		return false, erro.ErrHTTPRequest
	}
	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("error get Credentials")
		return false, erro.ErrHTTPRequest
	}
	//fmt.Println("***** ",credentials," === ",cfg.Region)
	//--------------------------
	// Set timeout APIGW
	client := &http.Client{Timeout: time.Second * 29}
		
	//Prepare the POST payload, transform interface to json
	requestPayload, err := json.Marshal(inter)
	if err != nil {
		childLogger.Error().Err(err).Msg("error json.Marshal(inter)")
		return false, erro.ErrBadRequest
	}

	// create a request POST data
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestPayload))
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return false, erro.ErrHTTPRequest
	}
	// Setup Headers
	req.Header.Add("Content-Type", "application/json;charset=UTF-8");
	req.Header.Add("x_apigw_api_id", autentication_data.AppClientID);
	req.Header.Add("x_appClient", autentication_data.ApiKeyID);
	// Sign the payload
	hash := sha256.Sum256([]byte(requestPayload))
	hexHash := fmt.Sprintf("%x", hash)
	
	signer := v4.NewSigner()
	err = signer.SignHTTP(ctx, credentials, req, hexHash, "execute-api", cfg.Region, time.Now())
	if err != nil {
		childLogger.Error().Err(err).Msg("error signer with credentials")
		return false, erro.ErrHTTPRequest
	}

	// Call a POST REST call
	resp, err := client.Do(req)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Do Request")
		return false, erro.ErrHTTPRequest
	}
	
	childLogger.Debug().Interface("Headers :", resp.Header).Msg("----")
	childLogger.Debug().Int("StatusCode :", resp.StatusCode).Msg("----")

	//var msg_err map[string]interface{}
	//json.NewDecoder(resp.Body).Decode(&msg_err)
	//fmt.Printf("%s\n", msg_err)
	//----
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

		// Decode result using the Interface
		result := inter
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			childLogger.Error().Err(err).Msg("error no ErrUnmarshal")
			return false, erro.ErrUnmarshal
		}
	
		return result, nil
}