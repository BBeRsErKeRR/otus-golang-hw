package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func size(path string) (int64, error) {
	i, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	stat, err := i.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func TestCopy(t *testing.T) {
	t.Run("unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", "/dev/null", 0, 1)
		require.ErrorIs(t, err, ErrUnsupportedFile, "Error should be: %v, got: %v", ErrUnsupportedFile, err)
	})

	t.Run("copy to /dev/null", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/dev/null", 0, 1)
		require.Nil(t, err, "Should not be errors")
	})

	tests := []struct {
		name   string
		input  string
		output string
		offset int64
		limit  int64
	}{
		{
			name:   "empty offset",
			input:  "testdata/input.txt",
			output: "empty_offset.txt",
			offset: 0,
			limit:  1,
		},
		{
			name:   "empty limit",
			input:  "testdata/out_offset0_limit0.txt",
			output: "empty_limit.txt",
			offset: 10,
			limit:  0,
		},
		{
			name:   "with offset and limit",
			input:  "testdata/input.txt",
			output: "with_offset_and_limit.txt",
			offset: 10,
			limit:  100,
		},
		{
			name:   "large offset",
			input:  "testdata/out_offset6000_limit1000.txt",
			output: "large_offset.txt",
			offset: 500,
			limit:  1000,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.CreateTemp("/tmp", tc.output)
			check(err)
			defer os.Remove(f.Name())

			s, err := size(tc.input)
			check(err)

			err = Copy(tc.input, f.Name(), tc.offset, tc.limit)
			require.Nil(t, err, "Should not be errors")
			resStat, err := f.Stat()
			check(err)
			var expected int64
			switch l := tc.limit; {
			case l == 0:
				expected = s - tc.offset
			case s-tc.offset > l:
				expected = tc.limit
			default:
				expected = s - tc.offset
			}

			require.Equal(t, expected, resStat.Size(), "Result file not equal current")
		})
	}
}
