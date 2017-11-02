package api

// User holds the user data retrieved from the User API.
type User struct {
	ID            string
	Email         string
	Organizations []string
	Roles         []string
}
