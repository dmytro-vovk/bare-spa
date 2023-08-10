//go:build sysboard

package network

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/errors"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/helpers/testtools/exec"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/uci"
)

func GetWAN() (*types.WAN, error) {
	const errFormat = "failed to get WAN"

	ip, err := getIP()
	if err != nil {
		return nil, errors.Wrap(err, errFormat)
	}

	mask, err := getMask()
	if err != nil {
		return nil, errors.Wrap(err, errFormat)
	}

	gateway, err := getGateway()
	if err != nil {
		return nil, errors.Wrap(err, errFormat)
	}

	ex := uci.New()
	wan := types.WAN{
		DHCP:    types.NewDHCP(ex.Get("network.wan.proto")),
		IP:      ip,
		Mask:    mask,
		Gateway: gateway,
	}
	if err := ex.Err(); err != nil {
		return nil, errors.CommandErr.Wrap(err, errFormat)
	}

	return &wan, nil
}

func SetWAN(wan types.WAN) error {
	err := uci.SetConfig(map[string]string{
		"network.wan.proto":   wan.DHCP.Proto(),
		"network.wan.ipaddr":  wan.IP,
		"network.wan.netmask": wan.Mask,
		"network.wan.gateway": wan.Gateway,
	})
	return errors.CommandErr.Wrap(err, "failed to set WAN")
}

func GetDNS() (*types.DNS, error) {
	dns1, dns2, err := getDNS()
	if err != nil {
		return nil, err
	}

	dns := types.DNS{
		DNS1: dns1,
		DNS2: dns2,
	}

	return &dns, nil
}

func SetDNS(dns types.DNS) error {
	err := uci.Set("network.wan.dns", uci.ToList([]string{dns.DNS1, dns.DNS2}))
	return errors.CommandErr.Wrap(err, "failed to set DNS")
}

func getIP() (string, error) {
	result, err := exec.Shell(`/sbin/ifconfig eth0 | grep "inet addr:" | cut -d ":" -f2 | awk '{ print $1 }'`)
	return result, errors.CommandErr.Wrap(err, "failed to get IP")
}

func getMask() (string, error) {
	result, err := exec.Shell(`/sbin/ifconfig eth0 | grep "Mask:" | cut -d ":" -f4 | awk '{ print $1 }'`)
	return result, errors.CommandErr.Wrap(err, "failed to get Mask")
}

func getGateway() (string, error) {
	result, err := exec.Shell(`/bin/netstat -r | grep "default" | awk '{ print $2 }'`)
	return result, errors.CommandErr.Wrap(err, "failed to get Gateway")
}

func getDNS() (dns1, dns2 string, err error) {
	ex := uci.New()
	switch dns := uci.AsArray(ex.Get("network.wan.dns")); {
	case len(dns) >= 2:
		dns2 = dns[1]
		fallthrough
	case len(dns) == 1:
		dns1 = dns[0]
	}

	err = errors.CommandErr.Wrap(ex.Err(), "failed to get DNS")
	return
}
