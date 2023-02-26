package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

var i *int

// Test the function on different structures and other types.
type (
	BadValidateInt struct {
		Test int `validate:"min:aa"`
	}
	BadValidateLength struct {
		Test string `validate:"len:[\b]"`
	}
	BadValidateRegexp struct {
		Test string `validate:"regexp:[\b]"`
	}
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$|len:20|nospaces"`
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

	Product struct {
		Name        string `validate:"len:15"`
		Application App    `validate:"nested"`
	}

	NewRules struct {
		Odd      int    `validate:"odd"`
		Even     int    `validate:"even"`
		BigInt   int64  `validate:"min:0"`
		NoSpaces string `validate:"nospaces"`
	}
)

func TestValidate(t *testing.T) {
	badParseValidatorTests := []struct {
		name string
		in   interface{}
	}{
		{
			name: "bad digit in max/min",
			in: BadValidateInt{
				Test: 12,
			},
		},
		{
			name: "bad regexp pattern",
			in: BadValidateRegexp{
				Test: "asd@test.ru",
			},
		},
		{
			name: "bad len",
			in: BadValidateLength{
				Test: "sometext",
			},
		},
	}
	for _, tt := range badParseValidatorTests {
		name := fmt.Sprintf("case %v", tt.name)
		t.Run(name, func(t *testing.T) {
			ve := ValidationErrors{}
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Error(t, err)
			require.NotErrorIs(t, err, ve)
		})
	}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: NewRules{
				Odd:      1,
				Even:     2,
				BigInt:   333,
				NoSpaces: "GoodText",
			},
		},
		{
			in: NewRules{
				Odd:      2,
				Even:     2,
				NoSpaces: "GoodText",
			},
			expectedErr: errorValidateOdd,
		},
		{
			in: NewRules{
				Odd:      1,
				Even:     3,
				NoSpaces: "GoodText",
			},
			expectedErr: errorValidateEven,
		},
		{
			in: NewRules{
				Odd:      1,
				Even:     2,
				NoSpaces: "Bad Text",
			},
			expectedErr: errorValidateNoSpaces,
		},
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
				Body: "bad res",
			},
			expectedErr: errorValidateIn,
		},
		{
			in: Response{
				Code: 200,
				Body: "good res",
			},
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
		{
			in: User{
				ID:     "test",
				Name:   "Test",
				Age:    23,
				Email:  "test2@somemail.tu",
				Role:   UserRole("stuff"),
				Phones: []string{"999-8888-8"},
			},
			expectedErr: errorValidateUnsupportedValueType,
		},
		{
			in: Product{
				Name: "Test",
				Application: App{
					Version: "0.0.1",
				},
			},
		},
		{
			in:          10,
			expectedErr: errorUnsupportedType,
		},
		{
			in:          make(chan []int),
			expectedErr: errorUnsupportedType,
		},

		{
			in:          i,
			expectedErr: errorUnsupportedType,
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
		})
	}
}
