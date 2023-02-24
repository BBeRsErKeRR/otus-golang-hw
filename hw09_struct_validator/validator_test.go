package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	BadValidateInt struct {
		test int `validate:"min:a5|max:5u0"`
	}
	BadValidateRegexp struct {
		test string `validate:"regexp:$\asdg\g^"` //nolint:staticcheck,govet
	}
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
	badParseValidatorTests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "Case bad digit in max/min",
			in: BadValidateInt{
				test: 12,
			},
		},
		{
			name: "Case bad regexp pattern",
			in: BadValidateRegexp{
				test: "asd@test.ru",
			},
		},
	}
	for _, tt := range badParseValidatorTests {
		t.Run(tt.name, func(t *testing.T) {
			ve := ValidationErrors{}
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Error(t, err)
			require.NotErrorIs(t, err, ve)
			_ = tt
		})
	}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: Token{
				Header:    []byte("Header...."),
				Payload:   []byte("Payload...."),
				Signature: []byte("Signature...."),
			},
		},
		{
			in: Response{
				Code: 201,
				Body: "",
			},
			expectedErr: errorValidateIn,
		},
		{
			in: User{
				ID:     "test",
				Name:   "Test",
				Age:    10,
				Email:  "test2@somemail.tu",
				Role:   UserRole("stuff"),
				Phones: []string{"999-8888-88-8"},
			},
			expectedErr: errorValidateMin,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr, fmt.Sprintf("Expected '%v', but not found.", tt.expectedErr))
			} else {
				require.NoError(t, err)
			}
			_ = tt
		})
	}
}
