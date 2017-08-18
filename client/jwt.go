// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "jwt-signin": jwt Resource Client
//
// Command:
// $ goagen
// --design=github.com/JormungandrK/jwt-issuer/design
// --out=$(GOPATH)/src/github.com/JormungandrK/jwt-issuer
// --version=v1.2.0-dirty

package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SigninJWTPath computes a request path to the signin action of jwt.
func SigninJWTPath() string {

	return fmt.Sprintf("/jwt/signin")
}

// Signs in the user and generates JWT token
func (c *Client) SigninJWT(ctx context.Context, path string, payload *Credentials) (*http.Response, error) {
	req, err := c.NewSigninJWTRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSigninJWTRequest create the request corresponding to the signin action endpoint of the jwt resource.
func (c *Client) NewSigninJWTRequest(ctx context.Context, path string, payload *Credentials) (*http.Request, error) {
	var body bytes.Buffer
	err := c.Encoder.Encode(payload, &body, "*/*")
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}
