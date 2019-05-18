package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/Microkubes/jwt-issuer/api"
	"github.com/Microkubes/jwt-issuer/app"
	"github.com/Microkubes/jwt-issuer/app/test"
	"github.com/Microkubes/jwt-issuer/config"
	"github.com/Microkubes/microservice-tools/gateway"
	"github.com/keitaroinc/goa"
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
	config := &config.Config{
		Jwt: config.JWTConfig{
			ExpiryTime:    10000, // 10 seconds
			Issuer:        "Mock Issuer",
			SigningMethod: "RS512",
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

	service := goa.New("signin")

	controller := NewSigninController(service, &MockUserAPI{
		Handler: func(user, pass string) (*api.User, error) {
			return &api.User{
				Email:         "email@example.com",
				ID:            "000000000001",
				Organizations: []string{"org1", "org2"},
				Roles:         []string{"user"},
				Active:        true,
			}, nil
		},
	}, NewMockKeyStore(), config)
	email := "someuser@test.com"
	pass := "pass"
	scope := "api:read"
	test.SigninJWTCreated(t, context.Background(), service, controller, &app.Credentials{
		Email:    &email,
		Password: &pass,
		Scope:    &scope,
	})
}

func TestSigninJWTBadRequest(t *testing.T) {
	config := &config.Config{
		Jwt: config.JWTConfig{
			ExpiryTime:    10000, // 10 seconds
			Issuer:        "Mock Issuer",
			SigningMethod: "RS512",
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
	service := goa.New("signin")

	controller := NewSigninController(service, &MockUserAPI{
		Handler: func(email, pass string) (*api.User, error) {
			return nil, nil // User not found
		},
	}, NewMockKeyStore(), config)
	email := "someuser@test.com"
	pass := "pass"
	scope := "api:read"
	test.SigninJWTBadRequest(t, context.Background(), service, controller, &app.Credentials{
		Email:    &email,
		Password: &pass,
		Scope:    &scope,
	})
}

func TestSigninJWTInternalServerError(t *testing.T) {
	config := &config.Config{
		Jwt: config.JWTConfig{
			ExpiryTime:    10000, // 10 seconds
			Issuer:        "Mock Issuer",
			SigningMethod: "RS512",
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
	service := goa.New("signin")

	controller := NewSigninController(service, &MockUserAPI{
		Handler: func(user, pass string) (*api.User, error) {
			return nil, fmt.Errorf("Test Error :)") // Return an error to cause internal server error
		},
	}, NewMockKeyStore(), config)
	email := "someuser@test.com"
	pass := "pass"
	scope := "api:read"

	test.SigninJWTInternalServerError(t, context.Background(), service, controller, &app.Credentials{
		Email:    &email,
		Password: &pass,
		Scope:    &scope,
	})
}
