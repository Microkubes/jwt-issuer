package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("jwt-signin", func() {
	Title("JWT Sign in")
	Description("Sign in and generate JWT token with claims")
	Version("1.0")
	Scheme("http")
	Host("localhost:8080")
})

var _ = Resource("jwt", func() {
	BasePath("jwt")
	Description("Sign in")

	Action("signin", func() {
		Description("Signs in the user and generates JWT token")
		Routing(POST("/signin"))
		Response(BadRequest, ErrorMedia)
		Response(Created)
	})

})
