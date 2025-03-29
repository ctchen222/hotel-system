package pgtypes

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinFirstNameLength = 2
	MinLastNameLength  = 2
	MinPasswordLength  = 7
)

type PGUser struct {
	Id                string `db:"id,omitempty" json:"id,omitempty"`
	FirstName         string `db:"firstname" json:"firstname"`
	LastName          string `db:"lastname" json:"lastname"`
	Email             string `db:"email" json:"email"`
	EncryptedPassword string `db:"encrypted_password" json:"encrypted_password,omitempty"`
}

type UpdateUserParams struct {
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
}

type CreateUserParams struct {
	FirstName string `db:"firstname" json:"firstname"`
	LastName  string `db:"lastname" json:"lastname"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < MinFirstNameLength {
		errors["firstName"] = "firstName must be at least 2 characters long"
	}
	if len(params.LastName) < MinLastNameLength {
		errors["lastName"] = "lastName must be at least 2 characters long"
	}
	if !IsEmailValid(params.Email) {
		errors["email"] = "email is invalid"
	}
	if len(params.Password) < MinPasswordLength {
		errors["password"] = "password must be at least 7 characters long"
	}
	return errors
}

func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func NewUserFromParams(params CreateUserParams) (*PGUser, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &PGUser{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
