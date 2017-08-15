package store

import (
	"fmt"
	"io/ioutil"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type KeyStore interface {
	GetPrivateKey() (interface{}, error)
	GetPrivateKeyByName(keyName string) (interface{}, error)
}

type FileKeyStore struct {
	PrivateKey interface{}
	KeysMap    map[string]interface{}
}

func (fks *FileKeyStore) GetPrivateKey() (interface{}, error) {
	if fks.PrivateKey != nil {
		return fks.PrivateKey, nil
	}
	return nil, fmt.Errorf("No default key loaded")
}

func (fks *FileKeyStore) GetPrivateKeyByName(keyName string) (interface{}, error) {
	priv, ok := fks.KeysMap[keyName]
	if !ok {
		return nil, fmt.Errorf("no key with name %s loaded", keyName)
	}
	return priv, nil
}

func NewFileKeyStore(keyFiles map[string]string) (KeyStore, error) {
	keyStore := FileKeyStore{
		KeysMap: make(map[string]interface{}),
	}
	for keyName, keyFile := range keyFiles {
		keyBytes, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		privKey, err := jwtgo.ParseRSAPrivateKeyFromPEM(keyBytes)
		if err != nil {
			return nil, err
		}
		keyStore.KeysMap[keyName] = privKey
	}
	defaultKey, ok := keyStore.KeysMap["default"]
	if !ok {
		return nil, fmt.Errorf("no default key for signing client JWT tokens defined")
	}
	keyStore.PrivateKey = defaultKey
	return &keyStore, nil
}
