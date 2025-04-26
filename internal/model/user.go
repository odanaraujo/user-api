package model

import (
	"errors"
	"net/mail"
	"time"
)

const (
	CacheTTL = time.Hour * 24
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Age   uint64 `json:"age"`
	Email string `json:"email"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.CPF == "" {
		return errors.New("cpf is required")
	}
	if u.Age == 0 {
		return errors.New("age is required")
	}
	if !ValidateAge(u.Age) {
		return errors.New("age must be less than or equal to 100")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if ValidateEmail(u.Email) {
		return errors.New("email is invalid")
	}
	return nil
}

func ValidateAge(age uint64) bool {
	return age <= 100
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err != nil
}
