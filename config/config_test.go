package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := `{
    "jwt":{
      "issuer": "Jormungandr JWT Authority",
      "signingMethod": "RS512",
      "expiryTime": 30000
    },
    "keys": {
      "default": "./keys/rsa_default",
      "system": "./keys/rsa_system"
    },
    "microservice": {
      "name": "jwt-issuer",
      "port": 8080,
      "virtual_host": "jwt.auth.jormugandr.org",
      "hosts": ["localhost", "jwt.auth.jormugandr.org"],
      "weight": 10,
      "slots": 100
    },
    "services": {
      "user-microservice": "http://localhost:8080/users"
    }
  }`
	cnfFile, err := ioutil.TempFile("", "tmp-config")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(cnfFile.Name())

	cnfFile.WriteString(config)

	cnfFile.Sync()

	loadedCnf, err := LoadConfig(cnfFile.Name())

	if err != nil {
		t.Fatal(err)
	}

	if loadedCnf == nil {
		t.Fatal("Configuration was not read")
	}
}
