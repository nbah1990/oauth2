package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashed), err
}

func IsHashSame(str string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str)) == nil
}
