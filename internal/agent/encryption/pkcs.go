package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type Encryptor struct {
	key []byte
}

func (e *Encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	publicKeyBlock, _ := pem.Decode(e.key)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}
