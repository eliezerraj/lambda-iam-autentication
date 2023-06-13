package core

// Type for restApi Call
type RestApiData struct {
	Host		string
	Path		string
}

// Type for Autentication
type Autentication struct {
	ClientID		string	`json:"access_key_id,omitempty"`
	CLientSecret	string	`json:"secret_access_key_id,omitempty"`
	ApiKeyID		string  `json:"x_apigw_api_id,omitempty"`
	AppClientID		string  `json:"x_appClient,omitempty"`
	Bearer			string  `json:"token,omitempty"`
}

// Constructor
func NewAutentication(options ...func(*Autentication)) *Autentication {
	x := &Autentication{}
	for _, o := range options {
	  o(x)
	}
	return x
}

func WithClientID(clientId string) func(*Autentication) {
	return func(s *Autentication) {
	  s.ClientID = clientId
	}
}
func WithCLientSecret(cLientSecret string) func(*Autentication) {
	return func(s *Autentication) {
	  s.CLientSecret = cLientSecret
	}
}
func WithApiKeyID(apiKeyId string) func(*Autentication) {
	return func(s *Autentication) {
	  s.ApiKeyID = apiKeyId
	}
}
func WithBearer(bearer string) func(*Autentication) {
	return func(s *Autentication) {
	  s.Bearer = bearer
	}
}
func WithAppClientID(appClientID string) func(*Autentication) {
	return func(s *Autentication) {
	  s.AppClientID = appClientID
	}
}