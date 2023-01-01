package crypto_utils

import (
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"log"
)

func Salt_the_earth(password string) string {
	salt := []byte{238, 36, 223, 36, 101, 231, 93, 10, 243, 28}
	//derived_key
	dk, err := scrypt.Key([]byte(password), salt, 16, 8, 1, 32)
	if err != nil {
		log.Printf(err.Error())
	}
	return base64.StdEncoding.EncodeToString(dk)
}
