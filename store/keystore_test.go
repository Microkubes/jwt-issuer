package store

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetPrivateKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatal(err)
	}
	ks := &FileKeyStore{
		PrivateKey: key,
	}

	pk, err := ks.GetPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	if pk == nil {
		t.Fatal("Key not found, but expected a valid key.")
	}
}

func TestGetPrivateKeyByName(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatal(err)
	}
	syskey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatal(err)
	}

	ks := &FileKeyStore{
		KeysMap: map[string]interface{}{
			"syskey": syskey,
		},
		PrivateKey: key,
	}

	pk, err := ks.GetPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	if pk == nil {
		t.Fatal("Key not found, but expected a valid key.")
	}

	sk, err := ks.GetPrivateKeyByName("syskey")
	if err != nil {
		t.Fatal(err)
	}
	if sk == nil {
		t.Fatal("syskey not found, but expected a valid key.")
	}

	_, err = ks.GetPrivateKeyByName("nokey")

	if err == nil {
		t.Fatal("Not expected to find the key")
	}
}

func tempKeyFile() (string, func()) {
	tf, err := ioutil.TempFile("", "keyfile")
	if err != nil {
		panic(err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		panic(err)
	}
	privKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(tf, privKey)
	if err != nil {
		panic(err)
	}
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestNewFileKeyStore(t *testing.T) {
	defKey, rmdef := tempKeyFile()
	defer rmdef()
	syskey, rmsys := tempKeyFile()
	defer rmsys()

	ks, err := NewFileKeyStore(map[string]string{
		"default": defKey,
		"system":  syskey,
	})

	if err != nil {
		t.Fatal(err)
	}
	if ks == nil {
		t.Fatal("New KeyStore was expected")
	}

	_, err = ks.GetPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	_, err = ks.GetPrivateKeyByName("system")
	if err != nil {
		t.Fatal(err)
	}
}
