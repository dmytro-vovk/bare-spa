package types

// TODO Пеоеделать на CreateDeviceDefault(d DeviceT) *Device
// Список устройств
var DeviceListBase = map[DeviceType]Device{
	Omega1: {
		Name:       "Омега-1",
		Type:       Omega1,
		Interfaces: []Interface{UART, SPI}, // Типы поддерживаемых интерфейсов
		Address:    1,                      // Default address
		State:      StateActive,
		Output:     nil, // При добавлении модуля номера регистров выходов добавляем сюда
		Input:      nil, // При добавлении модуля номера регистров входов добавляем сюда
		ModulesNum: 1,   // Максимум 1 подключаемый модуль (triac-2, triac-4, и т.д.)
		Temperature: []ThermometerSettings{
			{
				Name:         "Резистивный #1",
				Type:         Resistive,
				RegisterAddr: 0x0000,
				Interval:     5,
			},
			{
				Name:         "Резистивный #2",
				Type:         Resistive,
				RegisterAddr: 0x0001,
				Interval:     5,
			},
			{
				Name:         "Инфракрасный",
				Type:         Infrared,
				RegisterAddr: 0x0002,
				Interval:     5,
			},
			{
				Name:         "Термопара-0",
				Type:         Thermocouple,
				RegisterAddr: 0x0003,
				Interval:     5,
			},
			{
				Name:         "Термопара-1",
				Type:         Thermocouple,
				RegisterAddr: 0x0004,
				Interval:     5,
			},
		},
		Motors: []MotorSettings{
			{
				Name:              "Двигатель #1",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x0014, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x0012, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0010, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5мм при 200 шагах на оборот => 1 мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
			{
				Name:              "Двигатель #2",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x0015, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x0013, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0011, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5мм при 200 шагах на оборот => 1 мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
		},
		// Triac-2, Triac-4, Triac-6, Relay-2, Relay-4, Relay-6
		Modules: []Module{
			Triac2{ThermometerSettings{Name: "Triac-2"}},
			Triac4{ThermometerSettings{Name: "Triac-4"}},
			Triac6{ThermometerSettings{Name: "Triac-6"}},
			Relay2{ThermometerSettings{Name: "Relay-2"}},
			Relay4{ThermometerSettings{Name: "Relay-4"}},
			Relay6{ThermometerSettings{Name: "Relay-6"}}},
	},
	Incubator: {
		Name:       "Инкубатор",
		Type:       Incubator,
		Interfaces: []Interface{RS485}, // Типы поддерживаемых интерфейсов RS485 этот тот же UART + 1 пин
		Address:    0x7B,               // Default address #123 == 0x7B
		State:      StateActive,
		Output: []GPIO{
			{
				Name:      "OUT-1",
				Direction: DirectionOutput,
				Register:  0x0001,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-2",
				Direction: DirectionOutput,
				Register:  0x0002,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-3",
				Direction: DirectionOutput,
				Register:  0x0003,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "FAN-1",
				Direction: DirectionOutput,
				Register:  0x0004,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "FAN-2",
				Direction: DirectionOutput,
				Register:  0x0005,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
		}, // При добавлении модуля номера регистрорв выходов добавляем сюда
		Input: []GPIO{
			{
				Name:      "End-1",
				Direction: DirectionInput,
				Register:  0x0001,
			},
			{
				Name:      "End-2",
				Direction: DirectionInput,
				Register:  0x0002,
			},
			{
				Name:      "End-3",
				Direction: DirectionInput,
				Register:  0x0003,
			},
		}, // При добавлении модуля номера регистров входов добавляем сюда
		ModulesNum: 0, // Максимум 1 подключаемый модуль (triac-2, triac-4, и т.д.)
		Temperature: []ThermometerSettings{
			{
				Name:         "TMP-1 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0004,
				Interval:     5,
			},
			{
				Name:         "TMP-2 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0005,
				Interval:     5,
			},
			{
				Name:         "TMP-3 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0006,
				Interval:     5,
			},
			{
				Name:         "TMP-4 (DTH-XX)",
				Type:         DTH,
				RegisterAddr: 0x0007,
				Interval:     5,
			},
		},
		Motors: []MotorSettings{
			{
				Name:              "Двигатель #1",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x000A, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x0008, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0006, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
			{
				Name:              "Двигатель #2",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x000B, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x0009, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0007, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
		},
		Modules: nil,
	},
	Incubator2: {
		Name:       "Инкубатор",
		Type:       Incubator2,
		Interfaces: []Interface{Radio, RS485}, // Типы поддерживаемых интерфейсов RS485 этот тот-же UART + 1 пин
		Address:    0x7B,                      // Default address #123 == 0x7B
		State:      StateActive,
		Output: []GPIO{
			{
				Name:      "OUT-1",
				Direction: DirectionOutput,
				Register:  0x0001,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-2",
				Direction: DirectionOutput,
				Register:  0x0002,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-3",
				Direction: DirectionOutput,
				Register:  0x0003,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-4",
				Direction: DirectionOutput,
				Register:  0x0004,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-5",
				Direction: DirectionOutput,
				Register:  0x0005,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-6",
				Direction: DirectionOutput,
				Register:  0x0006,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-7",
				Direction: DirectionOutput,
				Register:  0x0007,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-8",
				Direction: DirectionOutput,
				Register:  0x0008,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
		}, // При добавлении модуля номера регистров выходов добавляем сюда
		Input: []GPIO{
			{
				Name:      "In-1",
				Direction: DirectionInput,
				Register:  0x0001,
			},
			{
				Name:      "In-2",
				Direction: DirectionInput,
				Register:  0x0002,
			},
			{
				Name:      "In-3",
				Direction: DirectionInput,
				Register:  0x0003,
			},
			{
				Name:      "In-4",
				Direction: DirectionInput,
				Register:  0x0004,
			},
			{
				Name:      "In-5",
				Direction: DirectionInput,
				Register:  0x0005,
			},
		}, // При добавлении модуля номера регистров входов добавляем сюда
		ModulesNum: 0, // Максимум 1 подключаемый модуль (triac-2, triac-4, и т.д.)
		Temperature: []ThermometerSettings{
			{
				Name:         "TMP-1 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0006,
				Interval:     5,
			},
			{
				Name:         "TMP-2 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0007,
				Interval:     5,
			},
			{
				Name:         "TMP-3 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0008,
				Interval:     5,
			},
			{
				Name:         "TMP-4 (DTH-XX)",
				Type:         DS18B20,
				RegisterAddr: 0x0009,
				Interval:     5,
			},
		},
		Motors: []MotorSettings{
			{
				Name:              "Двигатель #1",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x000F, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x000C, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0009, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5 мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
			{
				Name:              "Двигатель #2",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x0010, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x000D, // SM-1-TargetPosition Register Address
				SettingsReg:       0x000A, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5 мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
			{
				Name:              "Двигатель #3",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x00011, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x000E,  // SM-1-TargetPosition Register Address
				SettingsReg:       0x000B,  // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5 мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
		},
		Modules: []Module{
			Triac2{ThermometerSettings{Name: "Triac-2"}},
			Triac4{ThermometerSettings{Name: "Triac-4"}},
			Triac6{ThermometerSettings{Name: "Triac-6"}},
			Relay2{ThermometerSettings{Name: "Relay-2"}},
			Relay4{ThermometerSettings{Name: "Relay-4"}},
			Relay6{ThermometerSettings{Name: "Relay-6"}},
			Temp1{ThermometerSettings{Name: "Temp-1"}},
			Temp2{ThermometerSettings{Name: "Temp-2"}},
		},
	},
	EthIO2: {
		Name:       "Eth-IO-2",
		Type:       EthIO2,
		Interfaces: []Interface{RS485}, // Типы поддерживаемых интерфейсов RS485 этот тот-же UART + 1 пин
		Address:    150,                // Default address #123 == 0x7B
		State:      StateActive,
		Output: []GPIO{
			{
				Name:      "OUT-1",
				Direction: DirectionOutput,
				Register:  0x0001,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-2",
				Direction: DirectionOutput,
				Register:  0x0002,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
		}, // При добавлении модуля номера регистров выходов добавляем сюда
		Input: []GPIO{
			{
				Name:      "Input-1",
				Direction: DirectionInput,
				Register:  0x0001,
			},
			{
				Name:      "Input-2",
				Direction: DirectionInput,
				Register:  0x0002,
			},
			{
				Name:      "Input-3",
				Direction: DirectionInput,
				Register:  0x0003,
			},
			{
				Name:      "Input-4",
				Direction: DirectionInput,
				Register:  0x0004,
			},
		}, // При добавлении модуля номера регистров входов добавляем сюда
		ModulesNum: 1, // Максимум 1 подключаемый модуль (triac-2, triac4, и т.д.)
		Temperature: []ThermometerSettings{
			{
				Name:         "TMP-1 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0005,
				Interval:     5,
			},
			{
				Name:         "TMP-2 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0006,
				Interval:     5,
			},
		},
		Motors: []MotorSettings{
			{
				Name:              "Двигатель #1",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x000A, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x0008, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0006, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5 мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
		},
		Modules: []Module{
			Triac2{ThermometerSettings{Name: "Triac-2"}},
			Triac4{ThermometerSettings{Name: "Triac-4"}},
			Triac6{ThermometerSettings{Name: "Triac-6"}},
			Relay2{ThermometerSettings{Name: "Relay-2"}},
			Relay4{ThermometerSettings{Name: "Relay-4"}},
			Relay6{ThermometerSettings{Name: "Relay-6"}},
			Temp1{ThermometerSettings{Name: "Temp-1"}},
			Temp2{ThermometerSettings{Name: "Temp-2"}},
		},
	},
	EthIO4: {
		Name:       "Eth-IO-4",
		Type:       EthIO4,
		Interfaces: []Interface{RS485}, // Типы поддерживаемых интерфейсов RS485 этот тот-же UART+1пин
		Address:    150,                // Default address #123 == 0x7B
		State:      StateActive,
		Output: []GPIO{
			{
				Name:      "OUT-1",
				Direction: DirectionOutput,
				Register:  0x0001,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-2",
				Direction: DirectionOutput,
				Register:  0x0002,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-3",
				Direction: DirectionOutput,
				Register:  0x0003,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-4",
				Direction: DirectionOutput,
				Register:  0x0004,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
		}, // При добавлении модуля номера регистров выходов добавляем сюда
		Input: []GPIO{
			{
				Name:      "Input-1",
				Direction: DirectionInput,
				Register:  0x0001,
			},
			{
				Name:      "Input-2",
				Direction: DirectionInput,
				Register:  0x0002,
			},
			{
				Name:      "Input-3",
				Direction: DirectionInput,
				Register:  0x0003,
			},
			{
				Name:      "Input-4",
				Direction: DirectionInput,
				Register:  0x0004,
			},
		}, // При добавлении модуля номера регистров входов добавляем сюда
		ModulesNum: 1, // Максимум 1 подключаемый модуль (triac-2, triac-4, и т.д.)
		Temperature: []ThermometerSettings{
			{
				Name:         "TMP-1 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0005,
				Interval:     5,
			},
			{
				Name:         "TMP-2 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0006,
				Interval:     5,
			},
			{
				Name:         "TMP-3 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0005,
				Interval:     5,
			},
			{
				Name:         "TMP-4 (DTH-XX)",
				Type:         DTH,
				RegisterAddr: 0x0006,
				Interval:     5,
			},
		},
		Motors: []MotorSettings{
			{
				Name:              "Двигатель #1",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x000A, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x0008, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0006, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
			{
				Name:              "Двигатель #2",
				Type:              MotorStepping,
				Mode:              MotorModeLinear,
				CurPositionReg:    0x000B, // SM-1-CurrPosition Register Address
				TargetPositionReg: 0x0009, // SM-1-TargetPosition Register Address
				SettingsReg:       0x0007, // SM-1-Settings Register Address
				CurPosition:       0,
				TargetPosition:    0,
				IsActive:          StateDisabled,
				Direction:         MotorDirectionForward,
				StepNumberPerUnit: 200, // Резьба 5 мм при 200 шагах на оборот => 1мм == 40 шагов.
				// Резьба 0,75 мм при 400 шагах на оборот => 533.3333 шага на мм
				StepDelay: 1, // 1 ms => 1000 шагов в секунду => 5 оборотов/сек (при 200 шагов на оборот)
			},
		},
		Modules: []Module{
			Triac2{ThermometerSettings{Name: "Triac-2"}},
			Triac4{ThermometerSettings{Name: "Triac-4"}},
			Triac6{ThermometerSettings{Name: "Triac-6"}},
			Relay2{ThermometerSettings{Name: "Relay-2"}},
			Relay4{ThermometerSettings{Name: "Relay-4"}},
			Relay6{ThermometerSettings{Name: "Relay-6"}},
			Temp1{ThermometerSettings{Name: "Temp-1"}},
			Temp2{ThermometerSettings{Name: "Temp-2"}},
		},
	},
	UPS12: {
		Name:       "UPS12",
		Type:       UPS12,
		Interfaces: []Interface{RS485}, // Типы поддерживаемых интерфейсов RS485 этот тот-же UART + 1 пин
		Address:    150,                // Default address #123 == 0x7B
		State:      StateActive,
		Output: []GPIO{
			{
				Name:      "OUT-1",
				Direction: DirectionOutput,
				Register:  0x0001,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
			{
				Name:      "OUT-2",
				Direction: DirectionOutput,
				Register:  0x0002,
				//DefaultOutputLevel: GPIOOutLevelFloating,
			},
		}, // При добавлении модуля номера регистров выходов добавляем сюда
		Input: []GPIO{
			{
				Name:      "Input-1",
				Direction: DirectionInput,
				Register:  0x0001,
			},
			{
				Name:      "Input-2",
				Direction: DirectionInput,
				Register:  0x0002,
			},
		}, // При добавлении модуля номера регистров входов добавляем сюда
		ModulesNum: 0, // Максимум 1 подключаемый модуль (triac-2, triac4, и т.д.)
		Temperature: []ThermometerSettings{
			{
				Name:         "TMP-1 (DS18B20)",
				Type:         DS18B20,
				RegisterAddr: 0x0005,
				Interval:     5,
			},
		},
	},
}
