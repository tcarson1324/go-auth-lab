package passwords

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	bytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	passwordHash := string(hash)
	return passwordHash, nil
}

func Verify(password, hash string) (bool, error) {
	passwordBytes := []byte(password)
	hashBytes := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}
