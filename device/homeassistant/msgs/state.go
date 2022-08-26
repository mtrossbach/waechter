package msgs

type StateResult struct {
	BaseResult
	Result []SensorState `json:"result"`
}

type SensorState struct {
	Attributes Attributes `json:"attributes"`
	EntityID   string     `json:"entity_id"`
	State      string     `json:"state"`
}

type Attributes struct {
	DeviceClass  string `json:"device_class"`
	FriendlyName string `json:"friendly_name"`
	MotionValid  bool   `json:"motion_valid"`
}
