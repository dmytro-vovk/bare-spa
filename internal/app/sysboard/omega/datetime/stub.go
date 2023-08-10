//go:build !sysboard

package datetime

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
	"time"
)

var _t types.Timestamp

func Get() (*types.Timestamp, error) {
	var zero types.Timestamp
	if _t == zero {
		_t = types.Timestamp(time.Now())
	} else {
		_t = types.Timestamp(time.Time(_t).Add(time.Second))
	}

	return &_t, nil
}

func Set(t types.Timestamp) error {
	_t = t
	return nil
}

var _ntp = types.NTP{
	Servers: []string{
		"0.lede.pool.ntp.org",
		"1.lede.pool.ntp.org",
		"2.lede.pool.ntp.org",
	},
}

func GetNTP() (*types.NTP, error) {
	return &_ntp, nil
}

func SetNTP(ntp types.NTP) error {
	_ntp = ntp
	return nil
}
