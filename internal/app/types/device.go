package types

import (
	"time"
)

type DeviceType string

const (
	Omega1     DeviceType = "Omega-1"
	Incubator  DeviceType = "Incubator"
	Incubator2 DeviceType = "Incubator-2"
	EthIO2     DeviceType = "Eth-IO-2"
	EthIO4     DeviceType = "Eth-IO-4"
	UPS12      DeviceType = "UPS12"
)

type Device struct {
	ID          int                   `json:"id"`
	Name        string                `json:"name"`
	Type        DeviceType            `json:"type"`        // Тип устройства Omega1, Incubator, etc
	Interfaces  []Interface           `json:"interfaces"`  // Доступные интерфейсы на девайсе
	Address     byte                  `json:"address"`     // Адрес устройства для работы с ним
	Online      OnlineState           `json:"online"`      // Система считает, что модуль онлайн, т.к. не превышен Interval + lastSeen
	Output      []GPIO                `json:"output"`      // Выходы устройства, если есть
	Input       []GPIO                `json:"input"`       // Входы устройства, если есть
	Temperature []ThermometerSettings `json:"temperature"` // Температурные датчики, если есть
	Motors      []MotorSettings       `json:"motors"`      // Двигателя, если есть
	Network     DeviceNetwork         `json:"network"`     // Сетевые настройки для девайсов типа
	State       DeviceState           `json:"state"`       // Включен или отключен. Можно снимать в ремонт и т.д., но настройки останутся
	LastSeen    time.Time             `json:"lastSeen"`
	Interval    int                   `json:"interval"`   // Максимальный интервал между опросами. Если превышен -- зачит устройство потеряно
	ModulesNum  int                   `json:"modulesNum"` // Количество модулей поддерживаемых одновременно
	Modules     []Module              `json:"modules"`    // Подключенный модуль добавляем сюда. Если поддерживает больше одного, то тоже тут
	Sent        int                   `json:"sent"`       // Количество отправленных пакетов
	Received    int                   `json:"received"`   // Количество полученных пакетов
	ErrorsNum   int                   `json:"errorsNum"`  // Количество пакетов с ошибками по этому девайсу
}

type DeviceState bool

const (
	StateDisabled DeviceState = false
	StateActive   DeviceState = true
)

type OnlineState bool

const (
	StateOffline OnlineState = false
	StateOnline  OnlineState = true
)
