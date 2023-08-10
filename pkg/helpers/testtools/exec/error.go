package exec

import (
	"fmt"
	"strings"
)

type CommandError struct {
	Cmd string
	Err error
}

func (e CommandError) Error() string {
	const errFormat = `execution "%s" failed`
	trimmed := strings.TrimSpace(e.Err.Error())
	if trimmed == "" {
		return fmt.Sprintf(errFormat, e.Cmd)
	}

	var lineFeed string
	if strings.Contains(trimmed, "\n") {
		lineFeed = "\n"
	}

	return fmt.Sprintf(errFormat+": %s", e.Cmd, lineFeed+trimmed)
}
