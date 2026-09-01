package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pwd string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(HashPassword), nil
}
