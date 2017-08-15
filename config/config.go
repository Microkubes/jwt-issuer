package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/JormungandrK/microservice-tools/gateway"
)

type JWTConfig struct {
	SigningMethod string `json:"signingMethod"`
	Issuer        string `json:"issuer"`
	ExpiryTime    int    `json:"expiryTime"`
}

type Config struct {
	Jwt          JWTConfig                  `json:"jwt"`
	Microservice gateway.MicroserviceConfig `json:"microservice"`
	Services     map[string]string          `json:"services"`
	Keys         map[string]string          `json:"keys"`
}

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
