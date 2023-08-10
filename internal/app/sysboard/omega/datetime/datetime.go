//go:build sysboard

package datetime

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/errors"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/helpers/testtools/exec"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/uci"
	"time"
)

var date = exec.Command("date")

func Get() (*types.Timestamp, error) {
	result, err := date.Exec("--utc")
	if err != nil {
		return nil, errors.CommandErr.Wrap(err, "failed to get datetime")
	}

	dt, err := time.Parse(time.UnixDate, result)
	if err != nil {
		return nil, errors.CommandErr.Wrap(err, "failed to parse datetime")
	}

	t := types.Timestamp(dt)
	return &t, nil
}

func Set(t types.Timestamp) error {
	const errFormat = "failed to set datetime"

	ex := uci.New()
	if enabled := uci.Itob(ex.Get("system.ntp.enabled")); enabled {
		return errors.CommandErr.Wrap(errors.New("NTP enabled"), errFormat)
	}

	_, err := date.Exec("--set", t.String())
	return errors.CommandErr.Wrap(err, errFormat)
}

func GetNTP() (*types.NTP, error) {
	ex := uci.New()
	ntp := types.NTP{
		Enabled: uci.Itob(ex.Get("system.ntp.enabled")),
		Servers: uci.AsArray(ex.Get("system.ntp.server")),
	}
	if err := ex.Err(); err != nil {
		return nil, errors.CommandErr.Wrap(err, "failed to get NTP servers")
	}

	return &ntp, nil
}

func SetNTP(ntp types.NTP) error {
	err := uci.SetConfig(map[string]string{
		"system.ntp.enabled": uci.Btoi(ntp.Enabled),
		"system.ntp.server":  uci.ToList(ntp.Servers),
	})
	return errors.CommandErr.Wrap(err, "failed to set NTP servers")
}
