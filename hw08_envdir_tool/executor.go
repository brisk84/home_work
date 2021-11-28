package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecuteCmd(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmde := ExecuteCmd(cmd[0], cmd[1:]...)
	cmde.Stdin = os.Stdin
	cmde.Stdout = os.Stdout
	cmde.Stderr = os.Stderr
	cmde.Env = os.Environ()

	for k, v := range env {
		if v.NeedRemove {
			var newEnv []string
			for _, oldElem := range cmde.Env {
				if !strings.Contains(oldElem, k) {
					newEnv = append(newEnv, oldElem)
				}
			}
			cmde.Env = newEnv
			continue
		}
		cmde.Env = append(cmde.Env, k+"="+v.Value)
	}

	if err := cmde.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	return cmde.ProcessState.ExitCode()
}
