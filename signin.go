package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/JormungandrK/jwt-issuer/api"
	"github.com/JormungandrK/jwt-issuer/app"
	"github.com/JormungandrK/jwt-issuer/config"
	"github.com/JormungandrK/jwt-issuer/store"
	"github.com/JormungandrK/microservice-security/jwt"

	"github.com/goadesign/goa"
	uuid "github.com/satori/go.uuid"
)

// SigninController implements the signin resource.
type SigninController struct {
	*goa.Controller
	api.UserAPI
	store.KeyStore
	*config.Config
}

// NewSigninController creates a signin controller.
func NewSigninController(service *goa.Service, userAPI api.UserAPI, keyStore store.KeyStore, config *config.Config) *SigninController {
	return &SigninController{
		Controller: service.NewController("SigninController"),
		UserAPI:    userAPI,
		KeyStore:   keyStore,
		Config:     config,
	}
}

// Signin runs the signin action.
func (c *SigninController) Signin(ctx *app.SigninJWTContext) error {
	// SigninController_Signin: start_implement

	// Put your logic here
	username := ctx.RequestData.Request.FormValue("username")
	password := ctx.RequestData.Request.FormValue("password")
	scope := ctx.RequestData.Request.FormValue("scope")

	if username == "" || password == "" {
		return ctx.BadRequest(goa.ErrBadRequest("credentials-required"))
	}

	user, err := c.UserAPI.FindUser(username, password)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if user == nil {
		return ctx.BadRequest(goa.ErrBadRequest("invalid-credentials"))
	}
	key, err := c.KeyStore.GetPrivateKey()
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	signedToken, err := c.signToken(*user, scope, key)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	ctx.ResponseData.Header().Add("Authorization", fmt.Sprintf("Bearer %s", signedToken))

	// SigninController_Signin: end_implement
	return ctx.Created()
}

func (c *SigninController) signToken(user api.User, scope string, key interface{}) (string, error) {
	claims := map[string]interface{}{}
	// standard claims
	claims["jti"] = uuid.NewV4().String()
	claims["iss"] = c.Config.Jwt.Issuer
	claims["exp"] = time.Now().Add(time.Duration(c.Config.Jwt.ExpiryTime) * time.Millisecond).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = 0
	claims["sub"] = user.ID

	// scope
	claims["scopes"] = scope

	// non-standard claims
	claims["userId"] = user.ID
	claims["username"] = user.Username
	claims["roles"] = strings.Join(user.Roles[:], ",")
	claims["organizations"] = strings.Join(user.Organizations[:], ",")

	return jwt.SignToken(claims, c.Config.Jwt.SigningMethod, key)
}
