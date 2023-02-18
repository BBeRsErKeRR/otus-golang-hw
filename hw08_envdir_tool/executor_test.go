package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func generateEnvVars() map[string]EnvValue {
	res := make(map[string]EnvValue)
	items := []struct {
		name       string
		value      string
		needRemove bool
	}{
		{name: "FOO", value: "BAR"},
		{name: "UNSET", value: "", needRemove: true},
		{name: "HOME", value: "", needRemove: true},
	}

	for _, i := range items {
		var envVar EnvValue
		envVar.Value = i.name
		envVar.NeedRemove = i.needRemove
		res[i.name] = envVar
	}
	return res
}

func TestRunCmd(t *testing.T) {
	envVars := generateEnvVars()

	tests := []struct {
		name     string
		argument string
		exitCode int
	}{
		{name: "valid system env variable", argument: "PATH", exitCode: 0},
		{name: "valid env", argument: "FOO", exitCode: 0},
		{name: "unset system env", argument: "HOME", exitCode: 1},
		{name: "invalid data", argument: "NOTFOUND", exitCode: 1},
		{name: "unset env", argument: "UNSET", exitCode: 1},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := []string{"bash", "testdata/check_env.sh", tc.argument}
			exitCode := RunCmd(args, envVars)
			require.Equal(t, exitCode, tc.exitCode)
		})
	}
}
