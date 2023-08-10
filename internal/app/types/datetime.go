package types

import (
	"encoding/json"
	"github.com/Sergii-Kirichok/pr/internal/app/errors"
	"time"
)

type Datetime struct {
	Timestamp `json:"timestamp"`
	NTP       NTP `json:"ntp"`
}

const layout = "2006-01-02 15:04:05"

type Timestamp time.Time

func (t Timestamp) String() string {
	return time.Time(t).Format(layout)
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var ts string
	if err := json.Unmarshal(data, &ts); err != nil {
		return errors.ParsingErr.Wrapf(err, "%T JSON unmarshalling", t)
	}

	dt, err := time.Parse(layout, ts)
	if err != nil {
		return errors.ParsingErr.Wrapf(err, "%T time parsing", t)
	}

	*t = Timestamp(dt)
	return nil
}

type NTP struct {
	Enabled bool     `json:"enabled"`
	Servers []string `json:"servers"`
}

func (ntp *NTP) Validate() {
	for i, s := range ntp.Servers {
		if s == "" {
			ntp.Servers = append(ntp.Servers[:i], ntp.Servers[i+1:]...)
		}
	}
}
