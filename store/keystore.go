package store

import (
	"fmt"
	"io/ioutil"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type KeyStore interface {
	GetPrivateKey() (interface{}, error)
}

type FileKeyStore struct {
	PrivateKey interface{}
}

func (fks *FileKeyStore) GetPrivateKey() (interface{}, error) {
	if fks.PrivateKey != nil {
		return fks.PrivateKey, nil
	}
	return nil, fmt.Errorf("No key loaded")
}

func NewFileKeyStore(keyFile string) (KeyStore, error) {
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	privKey, err := jwtgo.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}
	return &FileKeyStore{PrivateKey: privKey}, nil
}
