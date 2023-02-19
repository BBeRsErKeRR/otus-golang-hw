package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readEnvValue(dir, fileName string) (EnvValue, error) {
	var envValue EnvValue

	file, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return envValue, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return envValue, err
	}

	envValue.NeedRemove = fi.Size() == 0
	if !envValue.NeedRemove {
		reader := bufio.NewReader(file)
		b, _, err := reader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) {
			return envValue, err
		}
		envValue.Value += strings.TrimRight(strings.ReplaceAll(string(b), "\x00", "\n"), " \t")
	}

	return envValue, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	res := make(Environment)
	for _, de := range dirFiles {
		if !de.IsDir() && !strings.Contains(de.Name(), "=") {
			res[de.Name()], err = readEnvValue(dir, de.Name())
			if err != nil {
				return nil, err
			}
		}
	}

	return res, nil
}
