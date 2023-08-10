package types

type Alias struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	DeviceID int    `json:"deviceID"`
	IsMT7688 bool   `json:"isMT7688"`
	GPIO     byte   `json:"gpio"`
}
