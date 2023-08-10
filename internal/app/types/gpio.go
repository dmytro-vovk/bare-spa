package types

type DirectionBit bool

const (
	DirectionOutput DirectionBit = true
	DirectionInput  DirectionBit = false
)

func (b DirectionBit) String() string {
	if b {
		return "output"
	}

	return "input"
}

type InversionBit bool

const (
	Inverted    InversionBit = true
	NonInverted InversionBit = false
)

func (b InversionBit) String() string {
	if b {
		return "inverted"
	}

	return "non-inverted"
}

type LevelBit bool

const (
	LevelHigh LevelBit = true
	LevelLow  LevelBit = false
)

func (b LevelBit) String() string {
	if b {
		return "high"
	}

	return "low"
}

type ModeBits struct {
	FirstBit uint
	LastBit  uint
}

type GPIO struct {
	Name      string       `json:"name"`
	Pin       byte         `json:"pin"`      // Для модулей и плат расширения номер по порядку
	Register  uint32       `json:"register"` // Номер регистра для чтения/записи для в устройств
	Level     LevelBit     `json:"level"`
	Direction DirectionBit `json:"direction"`
	Inversion InversionBit `json:"inversion"`
	// TODO в будущем добавить режим работы выхода, обычный или pwm
}
