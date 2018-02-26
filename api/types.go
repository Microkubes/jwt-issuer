package api

// User holds the user data retrieved from the User API.
type User struct {
	ID            string   `json:"id"`
	Email         string   `json:"email,omitempty"`
	Organizations []string `json:"organizations,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	Active        bool     `json:"active"`
	Namespaces    []string `json:"namespaces"`
}
