package tuya

import (
	"fmt"
)

type AlarmSensor struct {
	TuyaDevice
	Alarm bool
}

func NewAlarmSensor() Device {
	return &AlarmSensor{
		TuyaDevice: TuyaDevice{
			deviceType: ALARM_SENSOR,
		},
		Alarm: false,
	}
}

func (a *AlarmSensor) String() string {
	return fmt.Sprintf(`{"ip":"%s"}`, a.IP())
}

func (a *AlarmSensor) Start() error {
	return fmt.Errorf("started")
}
