//go:generate goagen bootstrap -d github.com/JormungandrK/jwt-issuer/design

package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/JormungandrK/jwt-issuer/app"
	"github.com/JormungandrK/microservice-tools/gateway"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("jwt-signin")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "signin" controller
	c := NewSigninController(service)
	app.MountJWTController(service, c)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}

type JWTConfig struct {
	SigningMethod string
	Issuer        string
	ExpiryTime    int
	KeyFile       string
}

type Config struct {
	jwt          JWTConfig
	microservice gateway.MicroserviceConfig
}

func LoadConfig(confFile string) (*Config, error) {
	confBytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(confBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
