package app

import (
	"github.com/Sergii-Kirichok/pr/internal/app/ifaces"
)

type Application struct {
	board ifaces.Sysboard
	notifier
}

func New() *Application {
	return &Application{}
}

type notifier interface {
	Notify(method string, params interface{})
}

func (a *Application) SetNotifier(n notifier) {
	a.notifier = n
}

func (a *Application) notify(method string, params interface{}, err error) error {
	if err != nil {
		return err
	}

	a.Notify(method, params)
	return nil
}

// todo: but they observe even if we haven't connections
func (a *Application) LaunchObservers() {
	for _, observer := range []func(){
		a.observeDatetime,
		a.observeGPIO,
		a.observeTemperature,
	} {
		go observer()
	}
}
