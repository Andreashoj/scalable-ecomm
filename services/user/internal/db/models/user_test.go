package models

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name     string
		email    string
		password string
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
			got, err := NewUser(tt.args.name, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		Name      string
		Password  string
		Email     string
		Role      Role
		CreatedAt time.Time
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Password:  tt.fields.Password,
				Email:     tt.fields.Email,
				Role:      tt.fields.Role,
				CreatedAt: tt.fields.CreatedAt,
			}
			if got := u.ComparePassword(tt.args.password); got != tt.want {
				t.Errorf("ComparePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_GetID(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		Name      string
		Password  string
		Email     string
		Role      Role
		CreatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Password:  tt.fields.Password,
				Email:     tt.fields.Email,
				Role:      tt.fields.Role,
				CreatedAt: tt.fields.CreatedAt,
			}
			if got := u.GetID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsValid(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		Name      string
		Password  string
		Email     string
		Role      Role
		CreatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Password:  tt.fields.Password,
				Email:     tt.fields.Email,
				Role:      tt.fields.Role,
				CreatedAt: tt.fields.CreatedAt,
			}
			if err := u.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createHashedPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createHashedPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("createHashedPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createHashedPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
