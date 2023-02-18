package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("testdata envVars", func(t *testing.T) {
		checks := []struct {
			name       string
			value      string
			needRemove bool
		}{
			{
				name:       "BAR",
				value:      "bar",
				needRemove: false,
			},
			{
				name:       "EMPTY",
				value:      "",
				needRemove: false,
			},
			{
				name:       "HELLO",
				value:      "\"hello\"",
				needRemove: false,
			},
			{
				name:       "UNSET",
				value:      "",
				needRemove: true,
			},
		}

		res, err := ReadDir("testdata/env")
		require.NoError(t, err)

		for _, c := range checks {
			check, ok := res[c.name]
			require.True(t, ok, "Could not get value %v from res", c.name)
			require.Equal(t, check.Value, c.value, "For %v result should be : '%v', got '%v'", c.name, c.value, check.Value)
			require.Equal(t, check.NeedRemove, c.needRemove, "For %v result should be : '%v'", c.name, c.needRemove)
		}
	})
}
