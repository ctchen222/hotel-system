package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinFirstNameLength = 2
	MinLastNameLength  = 2
	MinPasswordLength  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < MinFirstNameLength {
		errors["firstName"] = fmt.Sprintf("firstName must be at least %d characters long", MinFirstNameLength)
	}
	if len(params.LastName) < MinLastNameLength {
		errors["lastName"] = fmt.Sprintf("lastName must be at least %d characters long", MinLastNameLength)
	}
	if len(params.Password) < MinPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters long", MinPasswordLength)
	}
	if !IsEmailValid(params.Email) {
		errors["email"] = "email is invalid"
	}
	return errors
}

func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(encpw, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(password)) == nil
}

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

type UserUpdateParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
