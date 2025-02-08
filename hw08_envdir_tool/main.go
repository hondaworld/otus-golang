package main

import (
	"log"
	"os"
)

func main() {
	envDir := os.Args[1]

	environment, err := ReadDir(envDir)

	if err == nil {
		RunCmd(os.Args, environment)
	} else {
		log.Fatalf("failed to read files: %v", err)
	}
}
