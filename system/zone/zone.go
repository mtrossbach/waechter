package zone

import "github.com/mtrossbach/waechter/internal/config"

type Id string

type Zone struct {
	Id          Id     `json:"id"`
	DisplayName string `json:"displayName"`
	Perimeter   bool   `json:"perimeter"`
	Delayed     bool   `json:"delayed"`
	Armed       bool   `json:"-"`
}

func ZoneFromConfig(config config.ZoneConfig) Zone {
	return Zone{
		Id:          Id(config.Id),
		DisplayName: config.DisplayName,
		Perimeter:   config.Perimeter,
		Delayed:     config.Delayed,
		Armed:       false,
	}
}

const (
	NoZone Id = ""
)

func SubstitutionZone(displayName string, armed bool) Zone {
	return Zone{
		Id:          "_",
		DisplayName: displayName,
		Perimeter:   false,
		Delayed:     false,
		Armed:       armed,
	}
}
