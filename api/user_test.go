package api

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"reflect"
	"testing"

	"github.com/Microkubes/jwt-issuer/pkg/config"
	"github.com/Microkubes/microservice-tools/gateway"

	gock "gopkg.in/h2non/gock.v1"
)

type MockKeyStore struct {
	PrivateKey *rsa.PrivateKey
}

func (m *MockKeyStore) GetPrivateKey() (interface{}, error) {
	return m.PrivateKey, nil
}

func (m *MockKeyStore) GetPrivateKeyByName(keyName string) (interface{}, error) {
	return m.PrivateKey, nil
}

func NewMockStore() (*MockKeyStore, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &MockKeyStore{PrivateKey: privKey}, nil
}

func NewConfig() *config.Config {
	return &config.Config{
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
}

func TestSelfSignJWT(t *testing.T) {
	config := NewConfig()
	keyStore, err := NewMockStore()
	if err != nil {
		t.Fatal(err)
	}

	client := &UserAPIClient{
		UserServiceURL: "http://user.services.jormungandr.org:8001/user",
		KeyStore:       keyStore,
		Config:         config,
	}

	token, err := client.selfSignJWT()
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("Token was expected")
	}
	t.Log("Token: ", token)
}

func TestNewUserAPI(t *testing.T) {
	config := NewConfig()
	keyStore, err := NewMockStore()
	if err != nil {
		t.Fatal(err)
	}
	userAPI, err := NewUserAPI(config, keyStore)
	if err != nil {
		t.Fatal(err)
	}

	if userAPI == nil {
		t.Fatal("UserAPI was expected instead of nil")
	}

}

func TestFindUser(t *testing.T) {
	config := NewConfig()
	keyStore, err := NewMockStore()

	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{Transport: &http.Transport{}}

	userAPI := &UserAPIClient{
		UserServiceURL: "http://user.services.jormungandr.org:8001/user",
		KeyStore:       keyStore,
		Config:         config,
		Client:         client,
	}

	gock.New("http://user.services.jormungandr.org:8001").
		Post("/user/find").
		MatchType("json").
		JSON(map[string]string{"email": "some-mail", "password": "a-password"}).
		Reply(200).JSON(map[string]interface{}{
		"id":            "000000000001",
		"email":         "email@example.com",
		"organizations": []string{"org1", "org2"},
		"roles":         []string{"user", "admin"},
	})

	gock.InterceptClient(client)

	user, err := userAPI.FindUser("some-mail", "a-password")
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("Expected user")
	}
	if user.Email != "email@example.com" {
		t.Fatal("Email does not match expected")
	}
	if user.ID != "000000000001" {
		t.Fatal("ID does not match expected")
	}
	organizations := []string{"org1", "org2"}
	if !reflect.DeepEqual(user.Organizations, organizations) {
		t.Fatalf("Organizations do not match. Expected %v, but got %v.", organizations, user.Organizations)
	}
	roles := []string{"user", "admin"}
	if !reflect.DeepEqual(user.Roles, roles) {
		t.Fatalf("Roles do not match. Expected %v, but got %v.", roles, user.Roles)
	}
}
