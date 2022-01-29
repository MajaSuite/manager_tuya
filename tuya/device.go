package tuya

import (
	"errors"
	"fmt"
	"log"
)

const (
	NO_TYPE Type = iota
	ALARM_SENSOR
	BULB
	CEILING_LIGHT
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type Device interface {
	Type() Type
	Name() string
	IP() string
	SetIP(string)
	SendPacket() error
	String() string
	Start() error
}

type TuyaDevice struct {
	deviceType Type
	deviceName string
	deviceIp   string
}

type Type byte

func (t Type) String() string {
	switch t {
	case ALARM_SENSOR:
		return "Alarm sensor"
	case BULB:
		return "Bulb"
	case CEILING_LIGHT:
		return "Ceiling light"
	default:
		return "n/a"
	}
}

func NewDevice() Device {
	//switch expr {
	//case ALARM_SENSOR:
	return NewAlarmSensor()
	//default:
	return nil
	//}
}

func (td *TuyaDevice) Type() Type {
	return td.deviceType
}

func (td *TuyaDevice) Name() string {
	return td.deviceName
}

func (td *TuyaDevice) IP() string {
	return td.deviceIp
}

func (td *TuyaDevice) SetIP(ip string) {
	td.deviceIp = ip
	td.Start()
}

func (td *TuyaDevice) SendPacket() error {
	return ErrNotImplemented
}

func (td *TuyaDevice) String() string {
	return fmt.Sprintf(`{}`)
}

func (td *TuyaDevice) Start() error {
	log.Println("tuya start")
	return nil
}
