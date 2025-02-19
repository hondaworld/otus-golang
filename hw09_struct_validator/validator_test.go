package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:   "abcde4",
				Age:  18,
				Role: UserRole("admin"),
			},
			expectedErr: ErrInvalidLen,
		},
		{
			in: User{
				ID:   "abcde-abcde-abcde-abcde-abcde-abcdef",
				Age:  17,
				Role: UserRole("admin"),
			},
			expectedErr: ErrInvalidMin,
		},
		{
			in: User{
				ID:   "abcde-abcde-abcde-abcde-abcde-abcdef",
				Age:  55,
				Role: UserRole("admin"),
			},
			expectedErr: ErrInvalidMax,
		},
		{
			in: User{
				ID:    "abcde-abcde-abcde-abcde-abcde-abcdef",
				Age:   18,
				Role:  UserRole("admin"),
				Email: "john.doe@example.com",
			},
			expectedErr: ErrInvalidRegexp,
		},
		{
			in: User{
				ID:     "abcde-abcde-abcde-abcde-abcde-abcdef",
				Age:    18,
				Role:   UserRole("admin1"),
				Phones: []string{"+7999111111", "+7999111111"},
			},
			expectedErr: ErrInvalidIn,
		},
		{
			in: User{
				ID:     "abcde-abcde-abcde-abcde-abcde-abcdef",
				Age:    18,
				Role:   UserRole("admin"),
				Phones: []string{"+7999111111", "+79991"},
			},
			expectedErr: ErrInvalidLen,
		},
		{
			in: App{
				Version: "1.2",
			},
			expectedErr: ErrInvalidLen,
		},
		{
			in: Token{
				Header:    []byte{123, 2},
				Payload:   []byte{14, 44},
				Signature: []byte{55, 32},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
				Body: "12312",
			},
			expectedErr: ErrInvalidIn,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			var validationErrors ValidationErrors

			err := Validate(tt.in)

			if errors.As(err, &validationErrors) {
				for _, ve := range validationErrors {
					if !errors.Is(ve.Err, tt.expectedErr) {
						t.Errorf("expected error %v, got %v", tt.expectedErr, err)
					}
				}
			} else {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
