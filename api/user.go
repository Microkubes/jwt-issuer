package api

type UserAPI interface {
	FindUser(username, password string) (*User, error)
}
