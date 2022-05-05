package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Microkubes/jwt-issuer/api"
	"github.com/Microkubes/jwt-issuer/pkg/config"
	"github.com/Microkubes/jwt-issuer/pkg/store"
	"github.com/Microkubes/microservice-tools/gateway"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockKeyStore struct {
	key *rsa.PrivateKey
}

func (m *MockKeyStore) GetPrivateKey() (interface{}, error) {
	return m.key, nil
}

func (m *MockKeyStore) GetPrivateKeyByName(name string) (interface{}, error) {
	return m.key, nil
}

func NewMockKeyStore() *MockKeyStore {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return &MockKeyStore{
		key: key,
	}
}

type MockUserAPI struct {
	Handler func(email, password string) (*api.User, error)
}

func (m *MockUserAPI) FindUser(email, password string) (*api.User, error) {
	return m.Handler(email, password)
}

func TestSigninJWTCreated(t *testing.T) {
	conf := &config.Config{
		Jwt: config.JWTConfig{
			ExpiryTime:    10000, // 10 seconds
			Issuer:        "Mock Issuer",
			SigningMethod: "RS256",
		},
		Microservice: gateway.MicroserviceConfig{
			Hosts:            []string{"localhost", "jwt.auth.jormungandr.org"},
			MicroserviceName: "jwt-issuer",
			MicroservicePort: 8080,
			ServicesMaxSlots: 10,
			VirtualHost:      "jwt.auth.jormungandr.org",
			Weight:           10,
		},
		Services: map[string]string{
			"user-microservice": "http://user.services.jormungandr.org:8001/user",
		},
	}
	store.SetKeyStore(NewMockKeyStore())
	file, err := json.Marshal(conf)
	assert.NoError(t, err, "error marshaling conf data")
	err = ioutil.WriteFile("test.json", file, 0644)
	assert.NoError(t, err, "error writing test config data to file")
	defer os.RemoveAll("test.json")
	config.LoadConfig("test.json")
	api.SetUserApi(&MockUserAPI{
		Handler: func(user, pass string) (*api.User, error) {
			return &api.User{
				Email:         "email@example.com",
				ID:            "000000000001",
				Organizations: []string{"org1", "org2"},
				Roles:         []string{"user"},
				Active:        true,
			}, nil
		}},
	)
	test := echo.New()
	test.POST("/signin", Signin)
	form := url.Values{
		"email":    {"someuser@test.com"},
		"password": {"password"},
		"scope":    {"api:read"},
	}
	req, err := http.NewRequest("POST", "/signin", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err, "error creating request")
	rec := httptest.NewRecorder()
	c := test.NewContext(req, rec)
	if assert.NoError(t, Signin(c), "error in Signin handler") {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestSigninJWTBadRequest(t *testing.T) {
	conf := &config.Config{
		Jwt: config.JWTConfig{
			ExpiryTime:    10000, // 10 seconds
			Issuer:        "Mock Issuer",
			SigningMethod: "RS256",
		},
		Microservice: gateway.MicroserviceConfig{
			Hosts:            []string{"localhost", "jwt.auth.jormungandr.org"},
			MicroserviceName: "jwt-issuer",
			MicroservicePort: 8080,
			ServicesMaxSlots: 10,
			VirtualHost:      "jwt.auth.jormungandr.org",
			Weight:           10,
		},
		Services: map[string]string{
			"user-microservice": "http://user.services.jormungandr.org:8001/user",
		},
	}
	store.SetKeyStore(NewMockKeyStore())
	file, err := json.Marshal(conf)
	assert.NoError(t, err, "error marshaling conf data")
	err = ioutil.WriteFile("test.json", file, 0644)
	assert.NoError(t, err, "error writing test config data to file")
	defer os.RemoveAll("test.json")
	config.LoadConfig("test.json")
	api.SetUserApi(&MockUserAPI{
		Handler: func(user, pass string) (*api.User, error) {
			return nil, nil
		}},
	)
	test := echo.New()
	test.POST("/signin", Signin)
	form := url.Values{
		"email":    {"someuser@test.com"},
		"password": {"password"},
		"scope":    {"api:read"},
	}
	req, err := http.NewRequest("POST", "/signin", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err, "error creating request")
	rec := httptest.NewRecorder()
	c := test.NewContext(req, rec)
	if assert.NoError(t, Signin(c), "error in Signin handler") {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestSigninJWTUserNotActive(t *testing.T) {
	conf := &config.Config{
		Jwt: config.JWTConfig{
			ExpiryTime:    10000, // 10 seconds
			Issuer:        "Mock Issuer",
			SigningMethod: "RS256",
		},
		Microservice: gateway.MicroserviceConfig{
			Hosts:            []string{"localhost", "jwt.auth.jormungandr.org"},
			MicroserviceName: "jwt-issuer",
			MicroservicePort: 8080,
			ServicesMaxSlots: 10,
			VirtualHost:      "jwt.auth.jormungandr.org",
			Weight:           10,
		},
		Services: map[string]string{
			"user-microservice": "http://user.services.jormungandr.org:8001/user",
		},
	}
	store.SetKeyStore(NewMockKeyStore())
	file, err := json.Marshal(conf)
	assert.NoError(t, err, "error marshaling conf data")
	err = ioutil.WriteFile("test.json", file, 0644)
	assert.NoError(t, err, "error writing test config data to file")
	defer os.RemoveAll("test.json")
	config.LoadConfig("test.json")
	api.SetUserApi(&MockUserAPI{
		Handler: func(user, pass string) (*api.User, error) {
			return &api.User{
				Email:         "email@example.com",
				ID:            "000000000001",
				Organizations: []string{"org1", "org2"},
				Roles:         []string{"user"},
				Active:        false,
			}, nil
		}},
	)
	test := echo.New()
	test.POST("/signin", Signin)
	form := url.Values{
		"email":    {"someuser@test.com"},
		"password": {"password"},
		"scope":    {"api:read"},
	}
	req, err := http.NewRequest("POST", "/signin", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err, "error creating request")
	rec := httptest.NewRecorder()
	c := test.NewContext(req, rec)
	if assert.NoError(t, Signin(c), "error in Signin handler") {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
