package ifaces

import "github.com/Sergii-Kirichok/pr/internal/app/types"

type WiFiServer interface {
	GetWiFiAp() (*types.WiFiAp, error)
	SetWiFiAp(types.WiFiAp) error
}

type WiFiClient interface {
	GetWiFiCl() (*types.WiFiCl, error)
	SetWiFiCl(types.WiFiCl) error
}

type WAN interface {
	GetWAN() (*types.WAN, error)
	SetWAN(types.WAN) error
}

type DNS interface {
	GetDNS() (*types.DNS, error)
	SetDNS(types.DNS) error
}
