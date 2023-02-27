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
	UnsupportedValidateInt struct {
		Test  int `validate:"len:10"`
		Test2 int `validate:"undefined"`
	}
	BadValidateArgs struct {
		Test  int    `validate:"min:aa"`
		Test2 string `validate:"len:[\b]"`
		Test3 string `validate:"regexp:[\b]"`
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
	var i *int

	badParseValidatorTests := []struct {
		name string
		in   interface{}
	}{
		{name: "wrong validator", in: UnsupportedValidateInt{Test: 12}},
		{name: "bad validation args", in: BadValidateArgs{Test: 12, Test2: "sometext", Test3: "asd@test.ru"}},
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
		{in: NewRules{Odd: 1, Even: 2, BigInt: 333, NoSpaces: "GoodText"}},
		{in: NewRules{Odd: 2, Even: 2, NoSpaces: "GoodText"}, expectedErr: errorValidateOdd},
		{in: NewRules{Odd: 1, Even: 3, NoSpaces: "GoodText"}, expectedErr: errorValidateEven},
		{in: NewRules{Odd: 1, Even: 2, NoSpaces: "Bad Text"}, expectedErr: errorValidateNoSpaces},
		{
			in: Token{
				Header:    []byte("Header...."),
				Payload:   []byte("Payload...."),
				Signature: []byte("Signature...."),
			},
		},
		{in: Response{Code: 201, Body: "bad res"}, expectedErr: errorValidateIn},
		{in: Response{Code: 200, Body: "good res"}},
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
		{in: 10, expectedErr: errorUnsupportedType},
		{in: make(chan []int), expectedErr: errorUnsupportedType},
		{in: i, expectedErr: errorUnsupportedType},
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

	manyErrorsTests := []struct {
		in             interface{}
		expectedErrors []error
	}{
		{
			in: User{
				ID:     "BIIIIIGGGGGGGGGGGGGGGGGGGGGGLENGTHHHHHHHHHHHH",
				Name:   "Test",
				Age:    51,
				Email:  "test2somemail.tu",
				Role:   UserRole("stuff"),
				Phones: []string{"999-8888-8"},
			},
			expectedErrors: []error{
				errorValidateUnsupportedValueType,
				errorValidateStringLength,
				errorValidateMax,
				errorValidateStringRegexp,
			},
		},
	}

	for i, tt := range manyErrorsTests {
		t.Run(fmt.Sprintf("case many errors %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			for _, expected := range tt.expectedErrors {
				require.ErrorAs(t, err, &expected, fmt.Sprintf("Expected '%v', but not found.", expected))
			}
		})
	}
}
