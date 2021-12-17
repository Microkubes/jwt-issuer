//go:generate goagen bootstrap -d github.com/Microkubes/jwt-issuer/design

package main

import (
	"os"

	"github.com/Microkubes/jwt-issuer/api"
	"github.com/Microkubes/jwt-issuer/app"
	"github.com/Microkubes/jwt-issuer/config"
	"github.com/Microkubes/jwt-issuer/store"
	"github.com/Microkubes/microservice-tools/utils/healthcheck"
	"github.com/Microkubes/microservice-tools/utils/version"
	"github.com/keitaroinc/goa"
	"github.com/keitaroinc/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("jwt-signin")

	cf := os.Getenv("SERVICE_CONFIG_FILE")
	if cf == "" {
		cf = "/run/secrets/microservice_jwt_issuer_config.json"
	}
	config, err := config.LoadConfig(cf)
	if err != nil {
		service.LogError("config", "err", err)
		return
	}

	keyStore, err := store.NewFileKeyStore(config.Keys)
	if err != nil {
		panic(err)
	}

	userAPI, err := api.NewUserAPI(config, keyStore)
	if err != nil {
		panic(err)
	}

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	service.Use(healthcheck.NewCheckMiddleware("/healthcheck"))

	service.Use(version.NewVersionMiddleware(config.Version, "/version"))

	// Mount "signin" controller
	c := NewSigninController(service, userAPI, keyStore, config)
	app.MountJWTController(service, c)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}
