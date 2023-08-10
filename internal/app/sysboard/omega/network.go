package omega

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/sysboard/omega/network"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
	log "github.com/sirupsen/logrus"
)

func (s *Sysboard) GetWiFiAp() (*types.WiFiAp, error) {
	ap, err := network.GetWiFiAp()
	if err == nil {
		ap.TTL = s.Network.WiFiAp.TTL
		ap.DHCP = s.Network.WiFiAp.DHCP
	}

	return ap, err
}

func (s *Sysboard) SetWiFiAp(ap types.WiFiAp) error {
	if err := network.SetWiFiAp(ap); err != nil {
		return err
	}

	log.Println("New Wi-Fi access point configuration was set")
	s.Network.WiFiAp = ap
	return nil
}

func (s *Sysboard) GetWiFiCl() (*types.WiFiCl, error) {
	return network.GetWiFiCl()
}

func (s *Sysboard) SetWiFiCl(cl types.WiFiCl) error {
	if err := network.SetWiFiCl(cl); err != nil {
		return err
	}

	s.Network.WiFiCl = cl
	return nil
}

func (s *Sysboard) GetWAN() (*types.WAN, error) {
	return network.GetWAN()
}

func (s *Sysboard) SetWAN(wan types.WAN) error {
	if err := network.SetWAN(wan); err != nil {
		return err
	}

	log.Printf("New WAN configuration was set")
	s.Network.WAN = wan
	return nil
}

func (s *Sysboard) GetDNS() (*types.DNS, error) {
	return network.GetDNS()
}

func (s *Sysboard) SetDNS(dns types.DNS) error {
	if err := network.SetDNS(dns); err != nil {
		return err
	}

	log.Printf("New DNS configuration was set")
	s.Network.DNS = dns
	return nil
}
