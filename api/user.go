package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JormungandrK/jwt-issuer/config"
	"github.com/JormungandrK/jwt-issuer/store"
	"github.com/JormungandrK/microservice-security/jwt"
	"github.com/afex/hystrix-go/hystrix"
	uuid "github.com/satori/go.uuid"
)

type UserAPI interface {
	FindUser(username, password string) (*User, error)
}

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

type UserAPIClient struct {
	UserServiceURL string
	store.KeyStore
	*config.Config
	*http.Client
}

func (userAPI *UserAPIClient) FindUser(username, password string) (*User, error) {

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/find", userAPI.Config.Services["user-microservice"]), strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	token, err := userAPI.selfSignJWT()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	var resp *http.Response
	err = hystrix.Do("user-microservice.find-user", func() error {
		r, e := userAPI.Client.Do(req)
		if resp.StatusCode != 200 {
			return fmt.Errorf(resp.Status)
		}
		resp = r
		return e
	}, nil)

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	userResp := &map[string]interface{}{}

	err = json.NewDecoder(resp.Body).Decode(userResp)
	if err != nil {
		return nil, err
	}
	user := User{
		ID:       (*userResp)["id"].(string),
		Username: (*userResp)["username"].(string),
		Email:    (*userResp)["email"].(string),
	}
	roles, ok := (*userResp)["roles"]
	if ok {
		user.Roles = toStringArr(roles)
	}
	organizations, ok := (*userResp)["organizations"]
	if ok {
		user.Organizations = toStringArr(organizations)
	}
	return &user, nil
}

func toStringArr(val interface{}) []string {
	strArr, ok := val.([]string)
	if !ok {
		strArr = []string{}
	}
	return strArr
}

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