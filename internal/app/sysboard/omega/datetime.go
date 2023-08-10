package omega

import (
	"github.com/Sergii-Kirichok/pr/internal/app/sysboard/omega/datetime"
	"github.com/Sergii-Kirichok/pr/internal/app/types"
)

func (s *Sysboard) GetDatetime() (*types.Timestamp, error) {
	return datetime.Get()
}

func (s *Sysboard) SetDatetime(t types.Timestamp) error {
	if err := datetime.Set(t); err != nil {
		return err
	}

	s.Datetime.Timestamp = t
	return nil
}

func (s *Sysboard) GetNTP() (*types.NTP, error) {
	return datetime.GetNTP()
}

func (s *Sysboard) SetNTP(ntp types.NTP) error {
	if err := datetime.SetNTP(ntp); err != nil {
		return err
	}

	s.Datetime.NTP = ntp
	return nil
}
