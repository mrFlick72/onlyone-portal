// package security

// import (
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"testing"
// )

// func TestJwk_RsaPublicKey(t *testing.T) {
// 	client := new(MockedWebClientObject)

// 	body := `{
// 	  "keys": [ 
// 		{
// 		  "kty": "RSA",
// 		  "n": "o76AudS2rsCvlz_3D47sFkpuz3NJxgLbXr1cHdmbo9xOMttPMJI97f0rHiSl9stltMi87KIOEEVQWUgMLaWQNaIZThgI1seWDAGRw59AO5sctgM1wPVZYt40fj2Qw4KT7m4RLMsZV1M5NYyXSd1lAAywM4FT25N0RLhkm3u8Hehw2Szj_2lm-rmcbDXzvjeXkodOUszFiOqzqBIS0Bv3c2zj2sytnozaG7aXa14OiUMSwJb4gmBC7I0BjPv5T85CH88VOcFDV51sO9zPJaBQnNBRUWNLh1vQUbkmspIANTzj2sN62cTSoxRhSdnjZQ9E_jraKYEW5oizE9Dtow4EvQ",
// 		  "use": "sig",
// 		  "alg": "RS256",
// 		  "e": "AQAB",
// 		  "kid": "6a8ba5652a7044121d4fedac8f14d14c54e4895b"
// 		}
// 	  ]
// 	}
// 	`
// 	jwk := Jwk{
// 		Url:    "http://localhost/jwk",
// 		Client: client,
// 	}

// 	client.On("Get", &web.Request{
// 		Url: "http://localhost/jwk",
// 	}).Return(&web.Response{
// 		Body:   body,
// 		Status: 200,
// 	})

// 	_, err := jwk.RsaPublicKey()
// 	assert.Nil(t, err)

// }

// func TestJwk_JwkSets(t *testing.T) {
// 	client := new(MockedWebClientObject)

// 	body := `{
// 	  "keys": [ 
// 		{
// 		  "kty": "RSA",
// 		  "n": "o76AudS2rsCvlz_3D47sFkpuz3NJxgLbXr1cHdmbo9xOMttPMJI97f0rHiSl9stltMi87KIOEEVQWUgMLaWQNaIZThgI1seWDAGRw59AO5sctgM1wPVZYt40fj2Qw4KT7m4RLMsZV1M5NYyXSd1lAAywM4FT25N0RLhkm3u8Hehw2Szj_2lm-rmcbDXzvjeXkodOUszFiOqzqBIS0Bv3c2zj2sytnozaG7aXa14OiUMSwJb4gmBC7I0BjPv5T85CH88VOcFDV51sO9zPJaBQnNBRUWNLh1vQUbkmspIANTzj2sN62cTSoxRhSdnjZQ9E_jraKYEW5oizE9Dtow4EvQ",
// 		  "use": "sig",
// 		  "alg": "RS256",
// 		  "e": "AQAB",
// 		  "kid": "6a8ba5652a7044121d4fedac8f14d14c54e4895b"
// 		}
// 	  ]
// 	}
// 	`
// 	jwk := Jwk{
// 		Url:    "http://localhost/jwk",
// 	}

// 	client.On("Get", &web.Request{
// 		Url: "http://localhost/jwk",
// 	}).Return(&web.Response{
// 		Body:   body,
// 		Status: 200,
// 	})

// 	_, err := jwk.JwkSets()
// 	assert.Nil(t, err)
// }

// type MockedWebClientObject struct {
// 	mock.Mock
// }

// func (mock *MockedWebClientObject) Get(request *web.Request) (*web.Response, error) {
// 	called := mock.Called(request)
// 	return called.Get(0).(*web.Response), nil
// }