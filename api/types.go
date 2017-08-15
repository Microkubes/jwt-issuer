package api

type User struct {
	ID            string
	Username      string
	Email         string
	Organizations []string
	Roles         []string
}
