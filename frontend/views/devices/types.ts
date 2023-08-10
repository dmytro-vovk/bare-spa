export type DeviceTypeValue = `${DeviceType}`

export enum DeviceType {
	EthIO2     = "Eth-IO-2",
	EthIO4     = "Eth-IO-4",
	Incubator1 = "Incubator-1",
	Incubator2 = "Incubator-2",
	Omega1     = "Omega-1",
	UPS12      = "UPS12",
	Embedded   = "Embedded"
}

export type DeviceInterfaceValue = `${DeviceInterface}`

export enum DeviceInterface {
	I2C   = "I2C",
	RS485 = "RS485",
	Radio = "Radio",
	SPI   = "SPI",
	UART  = "UART",
}

export type Device = {
	readonly id: number
	readonly type: DeviceType
	name: string
	interface: DeviceInterface
	address: number
	modules: Array<Module>
};

// todo: use spec func
export type DeviceConfig = {
	readonly name: string
	readonly interface: DeviceInterface
	readonly address: number
	readonly modules?: Array<ModuleConfig>
}

export type ModuleTypeValue = `${ModuleType}`

export enum ModuleType {
	TypeDTH          = "DTH-XX",
	TypeDS18B20      = "DS18B20",
	TypeInfrared     = "Инфракрасный",
	TypeResistive    = "Резистивный",
	TypeThermocouple = "Термопара на MLX90614",
	TypeTriac2       = "Triac-2",
	TypeTriac4       = "Triac-4",
	TypeTriac6       = "Triac-6",
	TypeRelay2       = "Relay-2",
	TypeRelay4       = "Relay-4",
	TypeRelay6       = "Relay-6",
}

export type Module = {
	readonly id: number
	readonly type: ModuleType
	name: string
}

export type ModuleConfig = {
	readonly name: string,
}
