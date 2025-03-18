package api

import (
	"reflect"
	"testing"

	"github.com/ctchen1999/hotel-system/internal/db"
	"github.com/ctchen1999/hotel-system/internal/types"
	"github.com/gofiber/fiber/v2"
)

func TestNewAuthHandler(t *testing.T) {
	type args struct {
		userStore db.UserStore
	}
	tests := []struct {
		name string
		args args
		want *AuthHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthHandler(tt.args.userStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthHandler_HandleLogin(t *testing.T) {
	type fields struct {
		userStore db.UserStore
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthHandler{
				userStore: tt.fields.userStore,
			}
			if err := a.HandleLogin(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("AuthHandler.HandleLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	type args struct {
		user *types.User
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(tt.args.user); got != tt.want {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
