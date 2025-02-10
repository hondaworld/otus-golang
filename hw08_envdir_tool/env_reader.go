package main

import (
	"bytes"
	"fmt"
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

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environment := make(map[string]EnvValue)
	files, err := os.ReadDir(dir)

	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		content, err := os.ReadFile(filePath)

		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
		}

		lines := strings.Split(string(content), "\n")

		key := file.Name()
		result := ""

		if len(lines) > 0 {
			result = lines[0]
		}

		if result != "" {
			result = strings.TrimRight(result, " \t")
			result = string(bytes.Replace([]byte(result), []byte{0x00}, []byte{'\n'}, -1))

			environment[key] = EnvValue{
				Value:      result,
				NeedRemove: false,
			}
		} else {
			environment[key] = EnvValue{
				Value:      result,
				NeedRemove: true,
			}
		}
	}

	return environment, nil
}
