package app

import (
	"github.com/Sergii-Kirichok/pr/internal/app/types"
	log "github.com/sirupsen/logrus"
	"time"
)

func (a *Application) observeDatetime() {
	log.Println("Datetime observer was launched")
	for range time.NewTicker(time.Second).C {
		dt, err := a.board.GetDatetime()
		if err != nil {
			log.Printf("datetime observer error: %s", err)
			continue
		}

		a.Notify("datetime.set", struct {
			Timestamp types.Timestamp `json:"timestamp"`
		}{Timestamp: *dt})
	}
}

func (a *Application) GetDatetime() (*types.Timestamp, error) {
	return a.board.GetDatetime()
}

func (a *Application) SetDatetime(req struct {
	Timestamp types.Timestamp `json:"timestamp"`
}) error {
	return a.notify("datetime.set", req, a.board.SetDatetime(req.Timestamp))
}

func (a *Application) GetNTP() (*types.NTP, error) {
	return a.board.GetNTP()
}

func (a *Application) SetNTP(req struct {
	NTP types.NTP `json:"ntp"`
}) error {
	req.NTP.Validate()
	return a.notify("datetime.setNTP", req, a.board.SetNTP(req.NTP))
}
