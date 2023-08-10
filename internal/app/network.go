package app

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
)

func (a *Application) GetWiFiAp() (*types.WiFiAp, error) {
	return a.board.GetWiFiAp()
}

func (a *Application) SetWiFiAp(req types.WiFiAp) error {
	return a.notify("network.setWiFiAp", req, a.board.SetWiFiAp(req))
}

func (a *Application) GetWiFiCl() (*types.WiFiCl, error) {
	return a.board.GetWiFiCl()
}

func (a *Application) SetWiFiCl(req types.WiFiCl) error {
	return a.notify("network.setWiFiCl", req, a.board.SetWiFiCl(req))
}

func (a *Application) GetWAN() (*types.WAN, error) {
	return a.board.GetWAN()
}

func (a *Application) SetWAN(req types.WAN) error {
	if err := a.board.SetWAN(req); err != nil {
		return err
	}

	// todo: можем сразу не получить данные по dhcp (poller)
	wan, err := a.board.GetWAN()
	if err != nil {
		return err
	}

	req = *wan
	a.Notify("network.setWAN", req)
	return nil
}

func (a *Application) GetDNS() (*types.DNS, error) {
	return a.board.GetDNS()
}

func (a *Application) SetDNS(req types.DNS) error {
	return a.notify("network.setDNS", req, a.board.SetDNS(req))
}
