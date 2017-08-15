// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "jwt-signin": Application Controllers
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

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(goa.NewGobDecoder, "application/gob", "application/x-gob")
	service.Decoder.Register(goa.NewXMLDecoder, "application/xml")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// JWTController is the controller interface for the JWT actions.
type JWTController interface {
	goa.Muxer
	Signin(*SigninJWTContext) error
}

// MountJWTController "mounts" a JWT resource controller on the given service.
func MountJWTController(service *goa.Service, ctrl JWTController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSigninJWTContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Signin(rctx)
	}
	service.Mux.Handle("POST", "/jwt/signin", ctrl.MuxHandler("signin", h, nil))
	service.LogInfo("mount", "ctrl", "JWT", "action", "Signin", "route", "POST /jwt/signin")
}