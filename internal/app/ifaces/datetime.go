package ifaces

import "github.com/Sergii-Kirichok/pr/internal/app/types"

type Datetime interface {
	GetDatetime() (*types.Timestamp, error)
	SetDatetime(types.Timestamp) error
	GetNTP() (*types.NTP, error)
	SetNTP(types.NTP) error
}