package auth

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// GenerateConfirmationCode генерирует случайный код подтверждения
func GenerateConfirmationCode() (string, error) {
	bytes := make([]byte, 16) // 32 символа в шестнадцатеричном представлении
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword хэширует пароль с использованием bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
