package arm

type Mode string

const (
	Disarmed  Mode = "disarmed"
	Perimeter Mode = "armed-perimeter"
	All       Mode = "armed-all"
)
