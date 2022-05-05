//go:generate goagen bootstrap -d github.com/Microkubes/jwt-issuer/design

package main

import (
	"log"
	"os"

	"github.com/Microkubes/jwt-issuer/api"
	"github.com/Microkubes/jwt-issuer/handlers"
	"github.com/Microkubes/jwt-issuer/pkg/config"
	"github.com/Microkubes/jwt-issuer/pkg/store"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	e.POST("/signin", handlers.Signin)

	cf := os.Getenv("SERVICE_CONFIG_FILE")
	if cf == "" {
		cf = "/run/secrets/microservice_jwt_issuer_config.json"
	}
	config, err := config.LoadConfig(cf)
	if err != nil {
		log.Printf("error loading config %+v", err)
		return
	}

	keyStore, err := store.NewFileKeyStore(config.Keys)
	if err != nil {
		panic(err)
	}

	_, err = api.NewUserAPI(config, keyStore)
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(
		e.Start(":8080"),
	)
}
