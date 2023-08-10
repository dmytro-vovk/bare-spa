package exec

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

type command struct {
	name, as string
	args     []string
}

func Command(name string) *command {
	return &command{
		name: name,
		as:   name,
	}
}

func (c *command) As(name string) *command {
	c.as = name
	return c
}

func (c command) With(args ...string) *command {
	c.args = args
	return &c
}

func Shell(cmd string) (string, error) {
	return Command("sh").Exec("-c", cmd)
}

func (c command) Exec(args ...string) (string, error) {
	args = append(c.args, args...)
	cmd := execCommand(c.name, args...)
	str := c.str(args...)
	out, err := cmd.Output()
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			err = errors.New(string(e.Stderr))
		}

		return "", CommandError{Cmd: str, Err: err}
	}

	result := strings.TrimSpace(string(out))
	c.log(str, result)
	return result, nil
}

func (c command) str(args ...string) (cmd string) {
	return strings.Join(append([]string{c.as}, args...), " ")
}

func (c command) log(cmd string, out string) {
	if out == "" {
		out = "ok"
	}

	var lineFeed string
	if strings.Contains(out, "\n") {
		lineFeed = "\n"
	}

	log.Printf("%s -> %s", cmd, lineFeed+out)
}
