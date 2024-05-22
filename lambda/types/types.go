package types

import (
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStoreInsertUserDto struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(registerUser RegisterUserDto) (RegisterUserDto, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)

	if err != nil {
		return RegisterUserDto{}, err
	}

	return RegisterUserDto{
		Password: string(hashedPassword),
		Username: registerUser.Username,
	}, nil

}

func ValidatePassword(hashedPassword string, plainTextPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword)) == nil
}
