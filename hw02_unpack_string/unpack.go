package hw02unpackstring

import (
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

var strDigits = map[rune]int{'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9}

func GenerateRuneSlice(char rune, count int) []rune {
	// Generate n-time rune
	r := []rune{}
	for i := 1; i < count; i++ {
		r = append(r, char)
	}
	return r
}

func Unpack(input string) (string, error) {
	// Unpack string with slashes and digits

	r := []rune{}

	isSlashy := false
	isDigit := false

	var last rune

	for _, c := range input {
		switch {
			
		// Check slashy
		case c == '\\':
			isDigit = false
			if isSlashy {
				isSlashy = false
				r = append(r, c)
			} else {
				isSlashy = true
			}

		// Check any 0-9
		case c >= '0' && c <= '9':
			if isSlashy {
				r = append(r, c)
				isSlashy = false
			} else {
				count := strDigits[c]

				if last == 0 {
					return "", ErrInvalidString
				}

				if isDigit {
					return "", ErrInvalidString
				}

				if count > 0 {
					r = append(r, GenerateRuneSlice(last, count)...)
				} else {
					r = r[:len(r)-1]
				}

				isDigit = true
			}

		default:
			if isSlashy {
				return "", ErrInvalidString
			}

			isDigit = false
			r = append(r, c)
		}

		last = c
	}

	if isSlashy {
		return "", ErrInvalidString
	}

	return string(r), nil
}
