package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JormungandrK/jwt-issuer/config"
	"github.com/JormungandrK/jwt-issuer/store"
	"github.com/JormungandrK/microservice-security/jwt"
	"github.com/afex/hystrix-go/hystrix"
	uuid "github.com/satori/go.uuid"
)

// UserAPI defines operations for interacting with the User service API.
type UserAPI interface {
	// FindUser performs a lookup for a user by its email and password.
	FindUser(email, password string) (*User, error)
}

// NewUserAPI creates a UserAPI with a given configuration and a store.KeyStore.
func NewUserAPI(config *config.Config, keyStore store.KeyStore) (UserAPI, error) {
	serviceURL, ok := config.Services["user-microservice"]
	if !ok {
		return nil, fmt.Errorf("no URL for the User Microservice")
	}
	client := &http.Client{}
	return &UserAPIClient{
		UserServiceURL: serviceURL,
		KeyStore:       keyStore,
		Config:         config,
		Client:         client,
	}, nil
}

// UserAPIClient holds the data for the user microservice client
type UserAPIClient struct {
	// UserServiceURL is the base URL of the user microservice. This should be the API gateway exposed URL.
	UserServiceURL string

	// KeyStore is a reference to the store.KeyStore used for loading private keys
	store.KeyStore

	// Config is the microservice configuration
	*config.Config

	// Client is a HTTP clien implementation used for access to the user microservice
	*http.Client
}

// FindUser looks up a user by its email and password by calling the ```find``` action of the user microservice.
func (userAPI *UserAPIClient) FindUser(email, password string) (*User, error) {

	credentials := map[string]string{
		"email":    email,
		"password": password,
	}

	payload, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/find", userAPI.Config.Services["user-microservice"]), bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}

	token, err := userAPI.selfSignJWT()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	var resp *http.Response
	err = hystrix.Do("user-microservice.find-user", func() error {
		r, e := userAPI.Client.Do(req)
		resp = r
		if r.StatusCode != 200 {
			return fmt.Errorf(r.Status)
		}

		return e
	}, nil)

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if err != nil {
		println(err)
		return nil, err
	}

	userResp := &map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(userResp)
	if err != nil {
		return nil, err
	}
	user := User{
		ID:    (*userResp)["id"].(string),
		Email: (*userResp)["email"].(string),
	}
	roles, ok := (*userResp)["roles"]
	if ok {
		println("Roles: ", roles)
		user.Roles = toStringArr(roles)
	}
	organizations, ok := (*userResp)["organizations"]
	if ok {
		user.Organizations = toStringArr(organizations)
	}
	if active, ok := (*userResp)["active"]; ok {
		if active != nil {
			if _, ok := active.(bool); ok {
				user.Active = active.(bool)
			}
		}
	}
	return &user, nil
}

// toStringArr converts and array of interface{} to a string array
func toStringArr(val interface{}) []string {
	intfArr, ok := val.([]interface{})
	strArr := []string{}
	if !ok {
		return strArr
	}
	for _, intfV := range intfArr {
		if sval, ok := intfV.(string); ok {
			strArr = append(strArr, sval)
		}
	}
	return strArr
}

// selfSignJWT generates a JWT token which is self-signed with the system private key.
// This token is used for accesing the /user/find API on the user microservice.
func (userAPI *UserAPIClient) selfSignJWT() (string, error) {
	sysKey, err := userAPI.KeyStore.GetPrivateKeyByName("system")
	if err != nil {
		return "", err
	}
	signingMethod := userAPI.Config.Jwt.SigningMethod

	claims := map[string]interface{}{
		"iss":      userAPI.Config.Jwt.Issuer,
		"exp":      time.Now().Add(time.Duration(30) * time.Second).Unix(),
		"jti":      uuid.NewV4().String(),
		"nbf":      0,
		"sub":      "jwt-issuer",
		"scope":    "api:read",
		"userId":   "system",
		"username": "system",
		"roles":    "system",
	}

	sysToken, err := jwt.SignToken(claims, signingMethod, sysKey)

	return sysToken, err
}
