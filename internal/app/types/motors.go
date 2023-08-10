package types

type Motor string

const (
	MotorDC       Motor = "Двигатель постоянного тока"
	MotorStepping Motor = "Шаговый двигатель"
	MotorServo    Motor = "Сервопривод"
)

type MotorMode string

const (
	MotorModeLinear MotorMode = "Линейное перемещение, мм."
	MotorModeRadial MotorMode = "Радиальное вращение, град."
)

type MotorDirection bool

const (
	MotorDirectionForward  MotorDirection = true
	MotorDirectionBackward MotorDirection = false
)

type MotorSettings struct {
	Name              string         `json:"name"`
	Type              Motor          `json:"type"` // Тип мотора (шаговый, сервопривод)
	Mode              MotorMode      `json:"mode"` // (линейный, радиальный)
	CurPosition       uint16         `json:"curPosition"`
	CurPositionReg    uint16         `json:"curPositionReg"`
	TargetPosition    uint16         `json:"targetPosition"`
	TargetPositionReg uint16         `json:"targetPositionReg"`
	StepDelay         uint16         `json:"stepDelay"`   // Interval in ms
	SettingsReg       uint16         `json:"settingsReg"` // Ex. 0x80FF= Active, dir CW, step interval 255ms,
	IsActive          DeviceState    `json:"isActive"`
	StepNumberPerUnit uint16         `json:"stepNumberPerUnit"` // Количество шагов на 1 мм или на 1 градус
	Direction         MotorDirection `json:"direction"`         // Направление вращения
}
