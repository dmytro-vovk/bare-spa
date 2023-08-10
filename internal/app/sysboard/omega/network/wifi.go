//go:build sysboard

package network

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/errors"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/helpers/testtools/exec"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/uci"
	"strconv"
	"strings"
)

func GetWiFiAp() (*types.WiFiAp, error) {
	const errFormat = "failed to get Wi-Fi server"

	enabled, err := getApState()
	if err != nil {
		return nil, errors.Wrap(err, errFormat)
	}

	ex := uci.New()
	ap := types.WiFiAp{
		Enabled:    enabled,
		SSID:       ex.Get("wireless.ap.ssid"),
		Password:   ex.Get("wireless.ap.key"),
		Channel:    ex.Get("wireless.radio0.channel"),
		Encryption: types.WiFiEncryption(ex.Get("wireless.ap.encryption")),
		IP:         ex.Get("network.wlan.ipaddr"),
		Mask:       ex.Get("network.wlan.netmask"),
	}
	if err := ex.Err(); err != nil {
		return nil, errors.CommandErr.Wrap(err, errFormat)
	}

	return &ap, nil
}

func SetWiFiAp(wifi types.WiFiAp) error {
	err := uci.SetConfig(map[string]string{
		"wireless.ap.ssid":        wifi.SSID,
		"wireless.ap.key":         wifi.Password,
		"wireless.radio0.channel": strings.ToLower(wifi.Channel),
		"dhcp.wlan.leasetime":     strconv.FormatUint(uint64(wifi.TTL), 10) + "m",
		"wireless.ap.encryption":  string(wifi.Encryption),
		"network.wlan.ipaddr":     wifi.IP,
		"network.wlan.netmask":    wifi.Mask,
	})
	if switchErr := switchWiFi(wifi.Enabled, wifi.DHCP); err == nil {
		err = switchErr
	}
	return errors.CommandErr.Wrap(err, "failed to get Wi-Fi server")
}

func GetWiFiCl() (*types.WiFiCl, error) {
	return nil, errors.New("method not allowed")
}

func SetWiFiCl(types.WiFiCl) error {
	return errors.New("method not allowed")
}

func Setup(nw types.Network) error {
	err := uci.SetConfig(map[string]string{
		"wireless.radio0.country": "UA",
		"dhcp.wlan.start":         "1",
		"dhcp.wlan.limit":         "3",
		"dhcp.wlan.leasetime":     strconv.FormatUint(uint64(nw.WiFiAp.TTL), 10) + "m",
	})
	if switchErr := switchWiFi(nw.WiFiAp.Enabled, nw.WiFiAp.DHCP); err == nil {
		err = switchErr
	}
	return errors.CommandErr.Wrap(err, "network setup failed")
}

// getApState returns enabled/disabled (true/false) state for access point
func getApState() (bool, error) {
	const errFormat = "failed to get Wi-Fi access point state"

	ex := uci.New()
	apDisabled := uci.Itob(ex.Get("wireless.ap.disabled"))
	radio0Disabled := uci.Itob(ex.Get("wireless.radio0.disabled"))
	if err := ex.Err(); err != nil {
		return false, errors.CommandErr.Wrap(err, errFormat)
	}

	if apDisabled != radio0Disabled {
		return false, errors.Wrap(switchWiFi(false, types.DHCPOff), errFormat)
	}

	// or radio0Disabled, it doesn't matter because they are equal
	return apDisabled, nil
}

func switchWiFi(enabled bool, dhcp types.DHCP) error {
	const errFormat = "Wi-Fi switching failed"

	config := map[string]string{
		"wireless.ap.disabled":     uci.Btoi(!enabled),
		"wireless.radio0.disabled": uci.Btoi(!enabled),
	}

	if enabled && dhcp.IsOn() {
		err := uci.SetConfig(config)
		if dhcpErr := configureDHCP(types.DHCPOn); err == nil {
			err = dhcpErr
		}
		return errors.CommandErr.Wrap(err, errFormat)
	}

	_ = configureDHCP(types.DHCPOff) // may be an error, but it isn't serious
	return errors.CommandErr.Wrap(uci.SetConfig(config), errFormat)
}

// configureDHCP starts or stops DHCP server
func configureDHCP(dhcp types.DHCP) error {
	_, err := exec.Command("/etc/init.d/odhcpd").Exec(dhcp.Arg())
	return errors.CommandErr.Wrapf(err, "failed to %q DHCP", dhcp.Arg())
}
