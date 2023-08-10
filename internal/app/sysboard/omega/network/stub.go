//go:build !sysboard

package network

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
)

var (
	_ap  types.WiFiAp
	_cl  types.WiFiCl
	_wan types.WAN
	_dns types.DNS
)

func GetWiFiAp() (*types.WiFiAp, error) {
	return &_ap, nil
}

func SetWiFiAp(ap types.WiFiAp) error {
	_ap = ap
	return nil
}

func GetWiFiCl() (*types.WiFiCl, error) {
	return &_cl, nil
}

func SetWiFiCl(cl types.WiFiCl) error {
	_cl = cl
	return nil
}

func GetWAN() (*types.WAN, error) {
	return &_wan, nil
}

func SetWAN(wan types.WAN) error {
	_wan = wan
	return nil
}

func GetDNS() (*types.DNS, error) {
	return &_dns, nil
}

func SetDNS(dns types.DNS) error {
	_dns = dns
	return nil
}

func Setup(network types.Network) error {
	_ap = network.WiFiAp
	_cl = network.WiFiCl
	_wan = network.WAN
	_dns = network.DNS
	return nil
}
