package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const defaultFailedCode = 1

func filterEnvByKeys(env []string, envVars Environment) []string {
	res := make([]string, 0, len(env))
	res = append(res, env...)
	additional := make([]string, 0, len(envVars))
	for k, v := range envVars {
		if !v.NeedRemove {
			additional = append(additional, fmt.Sprintf("%v=%v", k, v.Value))
		} else {
			index := -1
			for i, envString := range env {
				if strings.HasPrefix(envString, fmt.Sprintf("%v=", k)) {
					index = i
					break
				}
			}
			if index >= 0 {
				res = append(res[:index], res[index+1:]...)
			}
		}
	}
	return append(res, additional...)
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	execCmd.Env = filterEnvByKeys(os.Environ(), env)
	err := execCmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok { //nolint
			return exitErr.ExitCode()
		}
		return defaultFailedCode
	}

	return
}
