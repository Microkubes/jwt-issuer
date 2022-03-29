package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Microkubes/jwt-issuer/api"
	"github.com/Microkubes/jwt-issuer/pkg/config"
	"github.com/Microkubes/jwt-issuer/pkg/store"
	"github.com/Microkubes/microservice-security/jwt"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func Signin(c echo.Context) error {
	u := api.User{}
	if err := c.Bind(&u); err != nil {
		log.Fatalf("error binding body %+v\n", err)
		return c.JSON(500, err)
	}
	user, err := api.GetUserApi().FindUser(u.Email, u.Password)
	if err != nil {
		log.Fatalf("error making a request to user api %+v\n", err)
		return c.JSON(500, err)
	}
	if user == nil {
		return c.JSON(400, "invalid-credentials")
	}
	if !user.Active {
		return c.JSON(400, "account-not-activated")
	}
	key, err := store.GetFileKeyStore().GetPrivateKey()
	if err != nil {
		return c.JSON(500, "error getting private key")
	}
	st, err := signToken(*user, u.Scope, key)
	if err != nil {
		log.Fatalf("error signing key %+v\n", err)
		return c.JSON(500, err)
	}
	bt := fmt.Sprintf("Bearer %s", st)
	c.Response().Header().Set("Authorization", bt)
	return c.JSON(201, echo.Map{
		"token": bt,
	})
}

func signToken(user api.User, scope string, key interface{}) (string, error) {
	claims := map[string]interface{}{}
	randUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	// standard claims
	claims["jti"] = randUUID.String()
	claims["iss"] = config.GetConfig().Jwt.Issuer
	claims["exp"] = time.Now().Add(time.Duration(config.GetConfig().Jwt.ExpiryTime) * time.Millisecond).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = 0
	claims["sub"] = user.ID

	// scope
	claims["scopes"] = scope

	// non-standard claims
	claims["userId"] = user.ID
	claims["username"] = user.Email
	claims["roles"] = strings.Join(user.Roles[:], ",")
	claims["organizations"] = strings.Join(user.Organizations[:], ",")
	if user.Namespaces != nil {
		claims["namespaces"] = strings.Join(user.Namespaces, ",")
	}

	// NOTE: get signing algorithm from config file
	return jwt.SignToken(claims, "RS256", key)
}
