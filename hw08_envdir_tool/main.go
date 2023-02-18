package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	emptyArgs = errors.New("empty arguments")
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(emptyArgs)
		os.Exit(1)
	}
	dirName := os.Args[1]
	cmdArgs := os.Args[2:]

	envVars, err := ReadDir(dirName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(RunCmd(cmdArgs, envVars))

}
