package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
	if dir[:len(dir)-1] != "/" {
		dir += "/"
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ret := make(Environment)
	for _, v := range files {
		if strings.Contains(v.Name(), "=") {
			return nil, fmt.Errorf("filename can't contain =")
		}

		file, err := os.Open(dir + v.Name())
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fileScanner := bufio.NewScanner(file)

		fileScanner.Scan()
		if err := fileScanner.Err(); err != nil {
			return nil, err
		}

		env := EnvValue{}
		env.Value = strings.TrimRight(fileScanner.Text(), " \t")
		env.Value = strings.ReplaceAll(env.Value, string(rune(0)), "\n")

		if v.Size() == 0 {
			env.NeedRemove = true
		}
		ret[v.Name()] = env
	}

	return ret, nil
}
