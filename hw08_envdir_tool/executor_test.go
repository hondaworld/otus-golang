package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("check environment", func(t *testing.T) {
		cmd := []string{"", "", "./testdata/echo.sh", "arg1=1", "arg2=2"}
		var environment Environment = map[string]EnvValue{}
		environment["BAR"] = EnvValue{Value: "bar", NeedRemove: false}

		code := RunCmd(cmd, environment)

		require.Equal(t, code, 0)
	})
	t.Run("without arguments", func(t *testing.T) {
		cmd := []string{"", "", "./testdata/echo.sh"}
		var environment Environment = map[string]EnvValue{}
		environment["BAR"] = EnvValue{Value: "bar", NeedRemove: false}

		code := RunCmd(cmd, environment)

		require.Equal(t, code, 0)
	})
}
