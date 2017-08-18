package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/JormungandrK/jwt-issuer/api"
	"github.com/JormungandrK/jwt-issuer/app"
	"github.com/JormungandrK/jwt-issuer/app/test"
	"github.com/JormungandrK/jwt-issuer/config"
	"github.com/JormungandrK/microservice-tools/gateway"
	"github.com/goadesign/goa"
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
	Handler func(username, password string) (*api.User, error)
}

func (m *MockUserAPI) FindUser(username, password string) (*api.User, error) {
	return m.Handler(username, password)
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
				Username:      user,
			}, nil
		},
	}, NewMockKeyStore(), config)
	user := "someuser"
	pass := "pass"
	scope := "api:read"
	test.SigninJWTCreated(t, context.Background(), service, controller, &app.Credentials{
		Username: &user,
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
		Handler: func(user, pass string) (*api.User, error) {
			return nil, nil // Username not found
		},
	}, NewMockKeyStore(), config)
	user := "someuser"
	pass := "pass"
	scope := "api:read"
	test.SigninJWTBadRequest(t, context.Background(), service, controller, &app.Credentials{
		Username: &user,
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
	user := "someuser"
	pass := "pass"
	scope := "api:read"

	test.SigninJWTInternalServerError(t, context.Background(), service, controller, &app.Credentials{
		Username: &user,
		Password: &pass,
		Scope:    &scope,
	})
}