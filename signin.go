package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Microkubes/jwt-issuer/api"
	"github.com/Microkubes/jwt-issuer/app"
	"github.com/Microkubes/jwt-issuer/config"
	"github.com/Microkubes/jwt-issuer/store"
	"github.com/Microkubes/microservice-security/jwt"

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
	email := ctx.Payload.Email //ctx.RequestData.Request.FormValue("email")
	password := ctx.Payload.Password
	scope := ctx.Payload.Scope //ctx.RequestData.Request.FormValue("scope")

	if email == nil || password == nil || scope == nil {
		return ctx.BadRequest(goa.ErrBadRequest("credentials-required: email, password, scope"))
	}

	if scope == nil {
		return ctx.BadRequest(goa.ErrBadRequest("scope-required"))
	}

	user, err := c.UserAPI.FindUser(*email, *password)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if user == nil {
		return ctx.BadRequest(goa.ErrBadRequest("invalid-credentials"))
	}

	if !user.Active {
		return ctx.BadRequest(goa.ErrBadRequest("account-not-activated"))
	}

	key, err := c.KeyStore.GetPrivateKey()
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	signedToken, err := c.signToken(*user, *scope, key)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	bearerToken := fmt.Sprintf("Bearer %s", signedToken)

	ctx.ResponseData.Header().Add("Authorization", bearerToken)

	// SigninController_Signin: end_implement
	return ctx.Created(bearerToken)
}

func (c *SigninController) signToken(user api.User, scope string, key interface{}) (string, error) {
	claims := map[string]interface{}{}
	randUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	// standard claims
	claims["jti"] = randUUID.String()
	claims["iss"] = c.Config.Jwt.Issuer
	claims["exp"] = time.Now().Add(time.Duration(c.Config.Jwt.ExpiryTime) * time.Millisecond).Unix()
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

	return jwt.SignToken(claims, c.Config.Jwt.SigningMethod, key)
}
