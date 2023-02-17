package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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
			require.NoError(t, err)
			defer os.Remove(f.Name())

			s, err := size(tc.input)
			require.NoError(t, err)

			err = Copy(tc.input, f.Name(), tc.offset, tc.limit)
			require.NoError(t, err)

			resStat, err := f.Stat()
			require.NoError(t, err)
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

	negativeTests := []struct {
		name       string
		input      string
		output     string
		offset     int64
		limit      int64
		e          error
		skipCreate bool
	}{
		{
			name:       "unsupported file",
			input:      "/dev/urandom",
			output:     "/dev/null",
			offset:     0,
			limit:      1,
			e:          ErrUnsupportedFile,
			skipCreate: true,
		},
		{
			name:   "negative offset",
			input:  "testdata/out_offset6000_limit1000.txt",
			output: "large_offset.txt",
			offset: -100,
			limit:  1000,
			e:      ErrorNegativeNumber,
		},
		{
			name:   "negative limit",
			input:  "testdata/out_offset6000_limit1000.txt",
			output: "large_offset.txt",
			offset: 0,
			limit:  -100,
			e:      ErrorNegativeNumber,
		},
	}

	for _, tc := range negativeTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var f *os.File
			var err error
			if tc.skipCreate {
				f, err = os.Open(tc.output)
				require.NoError(t, err)
			} else {
				f, err = os.CreateTemp("/tmp", tc.output)
				require.NoError(t, err)
			}
			defer os.Remove(f.Name())

			err = Copy(tc.input, f.Name(), tc.offset, tc.limit)
			require.ErrorIs(t, err, tc.e, "Error should be: %v, got: %v", tc.e, err)
		})
	}
}
