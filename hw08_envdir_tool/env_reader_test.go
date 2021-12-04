package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("testdata/env/")
	if err != nil {
		fmt.Println(err)
	}

	refEnv := make(Environment)
	refEnv["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
	refEnv["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
	refEnv["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
	refEnv["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
	refEnv["UNSET"] = EnvValue{Value: "", NeedRemove: true}

	require.Equal(t, refEnv, env)
}

func TestFileNameContainsEqual(t *testing.T) {
	f, _ := os.Create("testdata/env/test=123")
	f.Close()

	_, err := ReadDir("testdata/env/")
	if err != nil {
		fmt.Println(err)
	}
	os.Remove("testdata/env/test=123")

	require.Error(t, err)
}
