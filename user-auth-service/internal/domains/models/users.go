package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	StatusPending   string = "pending"
	StatusConfirmed string = "confirmed"
)

// User представляет модель пользователя
type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword проверяет введенный пароль с хэшем
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Validate проверяет корректность данных пользователя
func (u *User) Validate() error {
	if u.Username == "" || u.Email == "" || u.Password == "" {
		return errors.New("username, email и password обязательны")
	}
	// Здесь можно добавить дополнительную логику валидации, например, проверку формата email
	return nil
}
