package uci

import (
	"github.com/pkg/errors"
	"strings"
)

type executor struct {
	commits map[string]struct{}
	err     error
}

func New() *executor {
	return &executor{
		commits: make(map[string]struct{}),
	}
}

func (e *executor) Get(aim string) string {
	if e.err != nil {
		return ""
	}

	val, err := get(aim)
	if err != nil {
		e.err = errors.Wrapf(err, "uci: failed to get %q", aim)
		return ""
	}

	return val
}

func (e *executor) Set(aim, value string) *executor {
	if e.err != nil {
		return e
	}

	if out := e.Get(aim); out == value || e.err != nil {
		return e
	}

	if err := set(aim, value); err != nil {
		e.err = errors.Wrapf(err, "uci: failed to set %q=%q", aim, value)
		return e
	}

	e.commits[strings.Split(aim, ".")[0]] = struct{}{}
	return e
}

func Set(aim, value string) error {
	return New().Set(aim, value).Err()
}

func (e *executor) SetConfig(config map[string]string) *executor {
	for aim, value := range config {
		e.Set(aim, value)
	}

	return e
}

func SetConfig(config map[string]string) error {
	return New().SetConfig(config).Err()
}

func (e *executor) Err() error {
	if e.err != nil {
		return e.err
	}

	for config := range e.commits {
		if _, err := uci.Exec("commit", config); err != nil {
			return errors.Wrapf(err, "uci: failed to commit %q", config)
		}
	}

	if len(e.commits) != 0 {
		_, err := reloadConfig.Exec()
		return errors.Wrap(err, "uci: failed to reload config")
	}

	return nil
}

func ToList(array []string) string { return strings.Join(array, " ") }

func AsArray(value string) []string { return strings.Fields(value) }

func Itob(s string) bool { return s == "1" }

func Btoi(b bool) string {
	if b {
		return "1"
	}

	return "0"
}
