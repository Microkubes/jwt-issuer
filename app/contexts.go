// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "jwt-signin": Application Contexts
//
// Command:
// $ goagen
// --design=github.com/JormungandrK/jwt-issuer/design
// --out=$(GOPATH)/src/github.com/JormungandrK/jwt-issuer
// --version=v1.2.0-dirty

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// SigninJWTContext provides the jwt signin action context.
type SigninJWTContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Password *string
	Scope    *string
	Username *string
}

// NewSigninJWTContext parses the incoming request URL and body, performs validations and creates the
// context used by the jwt controller signin action.
func NewSigninJWTContext(ctx context.Context, r *http.Request, service *goa.Service) (*SigninJWTContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := SigninJWTContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramPassword := req.Params["password"]
	if len(paramPassword) > 0 {
		rawPassword := paramPassword[0]
		rctx.Password = &rawPassword
	}
	paramScope := req.Params["scope"]
	if len(paramScope) > 0 {
		rawScope := paramScope[0]
		rctx.Scope = &rawScope
	}
	paramUsername := req.Params["username"]
	if len(paramUsername) > 0 {
		rawUsername := paramUsername[0]
		rctx.Username = &rawUsername
	}
	return &rctx, err
}

// Created sends a HTTP response with status code 201.
func (ctx *SigninJWTContext) Created() error {
	ctx.ResponseData.WriteHeader(201)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *SigninJWTContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *SigninJWTContext) InternalServerError(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}