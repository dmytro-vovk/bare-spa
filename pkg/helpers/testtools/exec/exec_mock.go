//go:build !sysboard && exec

package exec

import (
	"os"
	"os/exec"
)

const (
	processTestEnv   = "GO_TEST_PROCESS"
	processTestMode  = "1"
	processDebugMode = "2"
)

func SetCommandMock(name string) func() {
	execCommand = commandMock(name)
	return func() { execCommand = exec.Command }
}

/*
commandMock is a function that initialises a new exec.Cmd which will
simply call Test<name>Process rather than the command it is provided. It will
also pass through the command and its arguments as an argument to Test<name>Process
*/
func commandMock(name string) func(command string, args ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=" + name, "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{processTestEnv + "=" + processTestMode}
		return cmd
	}
}

/*
TestProcessWrapper does all the dirty work for you, you need wrap EVERY function which emulates command execution.

Note: everything what you will print in your function will be interpreted as command's Stdout
Note: for testing your imitations you should manually set an environment variable processTestEnv to processDebugMode
*/
func TestProcessWrapper(exec func(cmd string, args []string) (success bool)) {
	env := os.Getenv(processTestEnv)
	switch env {
	case processTestMode:
	case processDebugMode:
	default:
		return
	}

	var (
		idx  = getCommandIndex(os.Args)
		cmd  = os.Args[idx]
		args = os.Args[idx+1:]
	)

	code := 1
	if exec(cmd, args) {
		code = 0
	}

	if env != processDebugMode || code == 1 {
		os.Exit(code)
	}
}

// getCommandIndex returns index from which REAL command starts
func getCommandIndex(args []string) int {
	for idx, arg := range args {
		if arg == "--" {
			return idx + 1
		}
	}

	return -1
}
