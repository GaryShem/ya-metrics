package encryption

import (
	"os"

	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

type EncryptorMiddleware interface {
	Encrypt(plaintext []byte) ([]byte, error)
}

var encryptor EncryptorMiddleware

func GetEncryptor() EncryptorMiddleware {
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
