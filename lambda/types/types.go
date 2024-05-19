package types

import (
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserWithHash struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type User struct {
	Username string `json:"username"`
}

func NewUser(registerUser RegisterUser) (UserWithHash, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)

	if err != nil {
		return UserWithHash{}, err
	}

	return UserWithHash{
		PasswordHash: string(hashedPassword),
		Username:     registerUser.Username,
	}, nil

}

func ValidatePassword(hashedPassword string, plainTextPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword)) == nil
}
