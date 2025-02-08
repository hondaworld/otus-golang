package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := cmd[2]
	args := cmd[3:]

	println(command)
	println(args)

	for key, envValue := range env {
		err := os.Unsetenv(key)
		if err != nil {
			log.Fatalf("failed to unset env var %s: %v", key, err)

			return 1
		}

		if !envValue.NeedRemove {
			err = os.Setenv(key, envValue.Value)
			if err != nil {
				log.Fatalf("failed to set env var %s: %v", key, err)

				return 1
			}
		}
	}

	com := exec.Command(command, args...)
	com.Stdout = os.Stdout
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin

	com.Env = os.Environ()

	err := com.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		log.Fatalf("failed to execute command: %v", err)

		return 1
	}

	return 0
}
