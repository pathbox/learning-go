package passhash

import (
	"golang.org/x/crypto/bcrypt"
)

func HashString(password string) (string, error) {
	key, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(key), nil
}

func HashBytes(password []byte) ([]byte, error) {
	key, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// MatchString returns true if the hash matches the password
func MatchString(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}

	return false
}

// MatchBytes returns true if the hash matches the password
func MatchBytes(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		return true
	}

	return false
}
