package device

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/api"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/model"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/msgs"
	"github.com/mtrossbach/waechter/system"
	"log"
)

type motionSensor struct {
	api           *api.Api
	entityId      string
	friendlyName  string
	systemControl system.Controller
	subId         uint64
}

func NewMotionSensor(api *api.Api, entityId string, friendlyName string) *motionSensor {
	return &motionSensor{api: api, entityId: entityId, friendlyName: friendlyName}
}

func (s *motionSensor) GetId() string {
	return s.entityId
}

func (s *motionSensor) GetDisplayName() string {
	return s.friendlyName
}

func (s *motionSensor) GetSubsystem() string {
	return model.SubsystemName
}

func (s *motionSensor) GetType() system.DeviceType {
	return system.MotionSensor
}

func (s *motionSensor) OnSystemStateChanged(state system.State, aMode system.ArmingMode, aType system.AlarmType) {

}

func (s *motionSensor) Setup(controller system.Controller) {
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
			s.systemControl.Alarm(system.BurglarAlarm, s)
		}
	}
}
