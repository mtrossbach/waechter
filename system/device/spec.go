package device

import "github.com/mtrossbach/waechter/internal/wslice"

type Spec struct {
	Id          Id
	DisplayName string
	Vendor      string
	Model       string
	Description string
	Sensors     []Sensor
	Actors      []Actor
}

func (s Spec) IsRelevant() bool {
	return len(s.Actors) > 0 || wslice.ContainsAny(s.Sensors,
		[]Sensor{MotionSensor, SmokeSensor, PanicSensor, TamperSensor, ArmingSensor, DisarmingSensor})
}
