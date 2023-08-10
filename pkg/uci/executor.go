//go:build sysboard || exec

package uci

import "github.com/Sergii-Kirichok/pr/pkg/helpers/testtools/exec"

var (
	uci          = exec.Command("uci")
	reloadConfig = exec.Command("reload_config")
)
