package device

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/homeassistant/api"
	"github.com/mtrossbach/waechter/device/homeassistant/msgs"
	"github.com/mtrossbach/waechter/system"
	"log"
)

type motionSensor struct {
	system.Device
	entityId      string
	api           *api.Api
	systemControl device.SystemController
	subId         uint64
}

func NewMotionSensor(device system.Device, api *api.Api, entityId string) *motionSensor {
	return &motionSensor{
		Device: system.Device{
			Id:   device.Id,
			Name: device.Name,
			Type: system.MotionSensor,
		},
		api:           api,
		systemControl: nil,
		entityId:      entityId,
		subId:         0,
	}
}

func (s *motionSensor) Setup(controller device.SystemController) {
	s.systemControl = controller
	go s.subscribe()
}

func (s *motionSensor) Teardown() {
	s.systemControl = nil

	if s.subId > 0 {
		s.api.UnsubscribeStateTrigger(s.subId)
		s.subId = 0
	}
}

func (s *motionSensor) subscribe() {
	_, c, id, _ := s.api.SubscribeStateTrigger(s.entityId)
	s.subId = id

	for data := range c {
		var event msgs.EventResponse
		err := json.Unmarshal(data, &event)
		if err != nil {
			log.Println(err.Error())
		}
		if event.Event.Variables.Trigger.ToState.State == "on" {
			s.systemControl.Alarm(system.BurglarAlarm, s.Device)
		}
	}
}
