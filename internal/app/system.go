package app

import (
	"github.com/Sergii-Kirichok/pr/internal/app/sysboard"
	"github.com/Sergii-Kirichok/pr/pkg/helpers/testtools/exec"
	"github.com/Sergii-Kirichok/pr/pkg/iface"
	log "github.com/sirupsen/logrus"
)

func (a *Application) GetInterfaces() ([]iface.Interface, error) {
	return a.board.GetInterfaces(), nil
}

func (a *Application) UseBoard(boardType string) *Application {
	var err error
	if a.board, err = sysboard.New(boardType); err != nil {
		log.Panicln(err)
	}

	if err = a.board.Setup(); err != nil {
		log.Panicln(err)
	}

	log.Printf("Using system board %s: %#v", boardType, a.board)
	return a
}

func (Application) LogRead() (string, error) {
	return exec.Shell(`logread | grep " webserver\[\d\+]: "`)
}
