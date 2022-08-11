package motionsensor

type statusPayload struct {
	Battery     int  `json:"battery"`
	Linkquality int  `json:"linkquality"`
	Occupancy   bool `json:"occupancy"`
	Tamper      bool `json:"tamper"`
}
