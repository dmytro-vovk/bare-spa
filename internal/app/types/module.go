package types

type SystemModule string

const (
	SysModuleRS485   SystemModule = "RS485"
	SysModuleRadio   SystemModule = "Radio"
	SysModuleModem3G SystemModule = "3G Modem"
)

type Module interface {
	Type() string
	// TODO
}

type ModuleSettings struct {
	Name        string                `json:"name"`
	Output      []GPIO                `json:"output"`
	Input       []GPIO                `json:"input"`
	Temperature []ThermometerSettings `json:"temperature"`
}

// Triacs

type Triac2 struct{ ThermometerSettings }

func (Triac2) Type() string { return `Triac-2` }

type Triac4 struct{ ThermometerSettings }

func (Triac4) Type() string { return `Triac-4` }

type Triac6 struct{ ThermometerSettings }

func (Triac6) Type() string { return `Triac-6` }

// Relays

type Relay2 struct{ ThermometerSettings }

func (Relay2) Type() string { return `Relay-2` }

type Relay4 struct{ ThermometerSettings }

func (Relay4) Type() string { return `Relay-4` }

type Relay6 struct{ ThermometerSettings }

func (Relay6) Type() string { return `Relay-6` }

// Temps

type Temp1 struct{ ThermometerSettings }

func (Temp1) Type() string { return `Temp-1` }

type Temp2 struct{ ThermometerSettings }

func (Temp2) Type() string { return `Temp-2` }
