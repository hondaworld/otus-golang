package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("check environment", func(t *testing.T) {
		dir := "./testdata/env"

		environment, err := ReadDir(dir)

		var result Environment = map[string]EnvValue{}
		result["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		result["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		result["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		result["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		result["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		require.Equal(t, result, environment)
		require.Nil(t, err)
	})
	t.Run("empty dir", func(t *testing.T) {
		dir := "./testdata/env1"

		environment, err := ReadDir(dir)

		require.Empty(t, environment)
		require.Nil(t, err)
	})
	t.Run("not exist dir", func(t *testing.T) {
		dir := "./testdata/env2"

		_, err := ReadDir(dir)

		require.NotNil(t, err)
	})
}
