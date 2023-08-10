package types

type WiFiEncryption string

const WPA WiFiEncryption = "psk"
const WPA2 WiFiEncryption = "psk2"

type DeviceNetwork struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type Network struct {
	WiFiAp WiFiAp `json:"wifiAp"`
	WiFiCl WiFiCl `json:"wifiCl"`
	WAN    WAN    `json:"wan"`
	DNS    DNS    `json:"dns"`
}

type WiFiAp struct {
	Enabled    bool           `json:"enabled"`
	SSID       string         `json:"ssid" validate:"wifi_ssid"`
	Password   string         `json:"password" validate:"password"`
	Channel    string         `json:"channel" validate:"wifi_channel"`
	TTL        uint           `json:"ttl" validate:"wifi_ttl"` // duration in minutes
	Encryption WiFiEncryption `json:"encryption" validate:"wifi_encryption"`
	IP         string         `json:"ip" validate:"required,ipv4"`
	Mask       string         `json:"mask" validate:"required,ipv4"`
	DHCP       DHCP           `json:"dhcp"`
}

type WiFiCl struct {
	Enabled  bool   `json:"enabled"`
	SSID     string `json:"ssid" validate:"required,min=1,max=32"`
	Password string `json:"password" validate:"required,alphanum,min=8,max=63"`
	IP       string `json:"ip" validate:"required,ipv4"`
	Mask     string `json:"mask" validate:"required,ipv4"`
	Gateway  string `json:"gateway" validate:"required,ipv4"`
	DHCP     DHCP   `json:"dhcp"`
}

type WAN struct {
	DHCP    DHCP   `json:"dhcp"`
	IP      string `json:"ip" validate:"required,ipv4"`
	Mask    string `json:"mask" validate:"required,ipv4"`
	Gateway string `json:"gateway" validate:"required,ipv4"`
}

type DNS struct {
	DNS1 string `json:"dns1" validate:"required,ipv4"`
	DNS2 string `json:"dns2,omitempty" validate:"omitempty,ipv4"`
}

type DHCP bool

const (
	DHCPOn  DHCP = true
	DHCPOff DHCP = false
)

func NewDHCP(proto string) DHCP { return proto == DHCPOn.Proto() }

func (d DHCP) IsOn() bool { return d == DHCPOn }

func (d DHCP) IsOff() bool { return d == DHCPOff }

func (d DHCP) Proto() string {
	if d {
		return "dhcp"
	}

	return "static"
}

func (d DHCP) Arg() string {
	if d {
		return "start"
	}

	return "stop"
}
