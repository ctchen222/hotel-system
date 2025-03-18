package types

import (
	"reflect"
	"testing"
)

func TestNewUserFromParams(t *testing.T) {
	type args struct {
		params CreateUserParams
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserFromParams(tt.args.params)
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
