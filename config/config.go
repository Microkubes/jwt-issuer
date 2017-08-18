package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/JormungandrK/microservice-tools/gateway"
)

// JWTConfig holds the JWT provider configuration.
type JWTConfig struct {
	// SigningMethod is the name of the method used for signing (RS256, RS384, RS512 etc)
	SigningMethod string `json:"signingMethod"`
	// Issuer is the name of the issuer (the name of the issuing server)
	Issuer string `json:"issuer"`
	// ExpiryTime sets the time period for which the JWT token is valid. The time is specified in milliseconds.
	ExpiryTime int `json:"expiryTime"`
}

// Config holds the microservice full configuration.
type Config struct {
	// Jwt is a JWTConfig for the JWT issuer.
	Jwt JWTConfig `json:"jwt"`
	// Microservice is a gateway.Microservice configuration for self-registration and service config.
	Microservice gateway.MicroserviceConfig `json:"microservice"`
	// Services is a map of <service-name>:<service base URL>. For example,
	// "user-microservice": "http://kong.gateway:8001/user"
	Services map[string]string `json:"services"`
	// Keys is a map <key name>:<key file path>. Should contain at least a "default" entry.
	// The keys are used for JWT generation and signing.
	// Example:
	// {
	//   "default": "keys/default.rsa.priv",
	//   "system": "keys/internal.system.key.rsa.priv"
	// }
	Keys map[string]string `json:"keys"`
}

// LoadConfig loads a Config from a configuration JSON file.
func LoadConfig(confFile string) (*Config, error) {
	confBytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(confBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
