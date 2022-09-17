package msgs

import "strings"

type StateResult struct {
	BaseResult
	Result []SensorState `json:"result"`
}

func (s *StateResult) GetEntityIdWithPrefixAndType(prefix string, stype string) string {
	if !strings.Contains(prefix, ".") {
		return ""
	}
	for _, s := range s.Result {
		if strings.HasPrefix(s.EntityID, prefix) && s.Attributes.DeviceClass == stype {
			return s.EntityID
		}
	}
	return ""
}

type SensorState struct {
	Attributes Attributes `json:"attributes"`
	EntityID   string     `json:"entity_id"`
	State      string     `json:"state"`
}

type Attributes struct {
	DeviceClass  string `json:"device_class"`
	FriendlyName string `json:"friendly_name"`
	MotionValid  *bool  `json:"motion_valid"`
}
