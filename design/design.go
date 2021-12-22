package design

import (
	. "github.com/keitaroinc/goa/design"
	. "github.com/keitaroinc/goa/design/apidsl"
)

var _ = API("jwt-signin", func() {
	Title("JWT Sign in")
	Description("Sign in and generate JWT token with claims")
	Version("1.0")
	Scheme("http")
	Host("localhost:8080")
	Consumes("application/x-www-form-urlencoded", func() {
		Package("github.com/keitaroinc/goa/encoding/form")
	})
})

var _ = Resource("jwt", func() {
	BasePath("/")
	Description("Sign in")

	Origin("*", func() {
		Methods("OPTIONS")
	})

	Action("signin", func() {
		Description("Signs in the user and generates JWT token")
		Payload(CredentialsPayload)
		Routing(POST("/signin"))
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
		Response(Created, String)
	})

})

// CredentialsPayload defines the credentials payload
var CredentialsPayload = Type("Credentials", func() {
	Attribute("email", String, "Credentials: email")
	Attribute("password", String, "Credentials: password")
	Attribute("scope", String, "Access scope (api:read, api:write)")
})
