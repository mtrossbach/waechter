package contactsensor

type statusPayload struct {
	Battery     int  `json:"battery"`
	Linkquality int  `json:"linkquality"`
	Contact     bool `json:"contact"`
	Tamper      bool `json:"tamper"`
}
