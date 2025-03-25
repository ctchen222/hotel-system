package types

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestIsValidPassword(t *testing.T) {
	type args struct {
		encpw    string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid Password",
			args: args{
				encpw:    "$2a$10$fXSf7.i3RluVG3GMGPa7FORF2NdWB9Els7veSo13teTYXChpVHJQG",
				password: "test1234",
			},
			want: true,
		},
		{
			name: "Invalid Password",
			args: args{
				encpw:    "1234",
				password: "test1234",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPassword(tt.args.encpw, tt.args.password); got != tt.want {
				t.Errorf("IsValidPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEmailValid(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid Email",
			args: args{
				email: "twobao@twobao.com",
			},
			want: true,
		},
		{
			name: "Invalid Email",
			args: args{
				email: "twobao@twobao",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmailValid(tt.args.email); got != tt.want {
				t.Errorf("isEmailValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserFromParams(t *testing.T) {
	pw := "TestPassword"
	encpw, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	type args struct {
		params CreateUserParams
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Valid User",
			args: args{
				params: CreateUserParams{
					FirstName: "TestFirstName",
					LastName:  "TestLastName",
					Email:     "EmailTest@gmail.com",
					Password:  pw,
				},
			},
			want: &User{
				FirstName:         "TestFirstName",
				LastName:          "TestLastName",
				Email:             "EmailTest@gmail.com",
				EncryptedPassword: string(encpw),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserFromParams(tt.args.params)
			got.EncryptedPassword = string(encpw)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserFromParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserFromParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateUserParams_Validate(t *testing.T) {
	type fields struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "Valid User",
			fields: fields{
				FirstName: "TestFirstName",
				LastName:  "TestLastName",
				Email:     "Twobao@gmail.com",
				Password:  "TestPassword",
			},
			want: map[string]string{},
		},
		{
			name: "InValid User",
			fields: fields{
				FirstName: "TestFirstName",
				LastName:  "",
				Email:     "Twobao@gmail.com",
				Password:  "TestPassword",
			},
			want: map[string]string{
				"lastName": fmt.Sprintf("lastName must be at least %d characters long", MinFirstNameLength),
			},
		},
		{
			name: "InValid User",
			fields: fields{
				FirstName: "TestFirstName",
				LastName:  "",
				Email:     "",
				Password:  "TestPassword",
			},
			want: map[string]string{
				"lastName": fmt.Sprintf("lastName must be at least %d characters long", MinLastNameLength),
				"email":    "email is invalid",
			},
		},
		{
			name: "InValid User",
			fields: fields{
				FirstName: "TestFirstName",
				LastName:  "",
				Email:     "",
				Password:  "",
			},
			want: map[string]string{
				"lastName": fmt.Sprintf("lastName must be at least %d characters long", MinLastNameLength),
				"email":    "email is invalid",
				"password": fmt.Sprintf("password must be at least %d characters long", MinPasswordLength),
			},
		},
		{
			name: "InValid User",
			fields: fields{
				FirstName: "",
				LastName:  "",
				Email:     "",
				Password:  "",
			},
			want: map[string]string{
				"firstName": fmt.Sprintf("firstName must be at least %d characters long", MinFirstNameLength),
				"lastName":  fmt.Sprintf("lastName must be at least %d characters long", MinLastNameLength),
				"email":     "email is invalid",
				"password":  fmt.Sprintf("password must be at least %d characters long", MinPasswordLength),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := CreateUserParams{
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
			}
			if got := params.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUserParams.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
