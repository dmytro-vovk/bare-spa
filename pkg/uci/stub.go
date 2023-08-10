//go:build !sysboard && !exec

package uci

import "github.com/Sergii-Kirichok/pr/pkg/helpers/testtools/exec"

var (
	uci          = exec.Command("true").As("uci")
	reloadConfig = exec.Command("true").As("reload_config")
)
