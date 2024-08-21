package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

var encryptor *Encryptor

func GetEncryptor() *Encryptor {
	return encryptor
}

func InitEncryptor(path string) error {
	logging.Log.Infoln("initializing encryption")
	if encryptor != nil {
		return nil
	}
	key, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	encryptor = &Encryptor{key}
	return nil
}

type Encryptor struct {
	key []byte
}

func (e *Encryptor) Encrypt(plain []byte) ([]byte, error) {
	publicKeyBlock, _ := pem.Decode(e.key)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plain)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}
