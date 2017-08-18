// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "jwt-signin": Application User Types
//
// Command:
// $ goagen
// --design=github.com/JormungandrK/jwt-issuer/design
// --out=$(GOPATH)/src/github.com/JormungandrK/jwt-issuer
// --version=v1.2.0-dirty

package client

// credentials user type.
type credentials struct {
	// Credentials: password
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	// Access scope (api:read, api:write)
	Scope *string `form:"scope,omitempty" json:"scope,omitempty" xml:"scope,omitempty"`
	// Credentials: username
	Username *string `form:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
}

// Publicize creates Credentials from credentials
func (ut *credentials) Publicize() *Credentials {
	var pub Credentials
	if ut.Password != nil {
		pub.Password = ut.Password
	}
	if ut.Scope != nil {
		pub.Scope = ut.Scope
	}
	if ut.Username != nil {
		pub.Username = ut.Username
	}
	return &pub
}

// Credentials user type.
type Credentials struct {
	// Credentials: password
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	// Access scope (api:read, api:write)
	Scope *string `form:"scope,omitempty" json:"scope,omitempty" xml:"scope,omitempty"`
	// Credentials: username
	Username *string `form:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
}
