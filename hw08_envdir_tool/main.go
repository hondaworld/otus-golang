package main

import (
	"fmt"
	"os"
)

func main() {
	envDir := os.Args[1]

	environment, err := ReadDir(envDir)

	if err != nil {
		fmt.Println(fmt.Errorf("failed to read files: %v", err))
		os.Exit(1)
	}

	code := RunCmd(os.Args, environment)
	os.Exit(code)
}
