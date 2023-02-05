package msgs

type EventResponse struct {
	BaseMessage
	Event EventContainer `json:"event"`
}

type EventContainer struct {
	Variables VariablesContainer `json:"variables"`
}

type VariablesContainer struct {
	Trigger TriggerContainer `json:"trigger"`
}

type TriggerContainer struct {
	Platform  string      `json:"platform"`
	EntityID  string      `json:"entity_id"`
	FromState SensorState `json:"from_state"`
	ToState   SensorState `json:"to_state"`
}
