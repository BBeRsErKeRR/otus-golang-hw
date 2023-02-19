package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const defaultFailedCode = 1

func removeInEnv(env []string, envKey string) []string {
	res := env
	index := -1
	for i, envString := range env {
		if strings.HasPrefix(envString, fmt.Sprintf("%v=", envKey)) {
			index = i
			break
		}
	}
	if index >= 0 {
		res = append(res[:index], res[index+1:]...)
	}
	return res
}

func filterEnvByKeys(env []string, envVars Environment) []string {
	res := make([]string, 0, len(envVars)+len(env))
	res = append(res, env...)
	for k, v := range envVars {
		res = removeInEnv(res, k)
		if !v.NeedRemove {
			res = append(res, fmt.Sprintf("%v=%v", k, v.Value))
		}
	}
	return res
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var execCmd *exec.Cmd

	switch length := len(cmd); length {
	case 0:
		fmt.Println("Not found execution command!")
		return 1
	case 1:
		execCmd = exec.Command(cmd[0]) //nolint:gosec
	default:
		execCmd = exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	}

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
