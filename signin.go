package main

import (
	"fmt"

	"github.com/JormungandrK/jwt-issuer/api"
	"github.com/JormungandrK/jwt-issuer/app"
	"github.com/JormungandrK/jwt-issuer/store"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
)

// SigninController implements the signin resource.
type SigninController struct {
	*goa.Controller
	api.UserAPI
	store.KeyStore
}

// NewSigninController creates a signin controller.
func NewSigninController(service *goa.Service) *SigninController {
	return &SigninController{Controller: service.NewController("SigninController")}
}

// Signin runs the signin action.
func (c *SigninController) Signin(ctx *app.SigninJWTContext) error {
	// SigninController_Signin: start_implement

	// Put your logic here
	username := ctx.RequestData.Request.FormValue("username")
	password := ctx.RequestData.Request.FormValue("password")
	if username == "" || password == "" {
		return ctx.BadRequest(fmt.Errorf("Credentials required"))
	}

	user, err := c.UserAPI.FindUser(username, password)
	if err != nil {

	}
	if user == nil {
		return ctx.BadRequest(fmt.Errorf("No user for credentials"))
	}

	// SigninController_Signin: end_implement
	return ctx.Created()
}

func (c *SigninController) signToken(claims *jwtgo.Claims) (string, error) {
	token := jwtgo.New(jwtgo.SigningMethodRS512)
}
