package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Model
	Firstname    string `json:"first_name"`
	Lastname     string `json:"last_name"`
	Email        string `json:"email"`
	Password     []byte `json:"-"`
	IsAmbassador bool   `json:"-"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
