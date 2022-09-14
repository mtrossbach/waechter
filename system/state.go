package system

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"os"
	"path"
)

type systemstate struct {
	State State      `json:"state"`
	Mode  ArmingMode `json:"mode"`
	Alarm AlarmType  `json:"alarm"`
}

func saveState(state State, mode ArmingMode, alarm AlarmType) {
	s := systemstate{
		State: state,
		Mode:  mode,
		Alarm: alarm,
	}

	data, err := json.Marshal(s)
	if err != nil {
		log.Error().Err(err).Msg("Could not save state")
		return
	}

	err = os.WriteFile(path.Join(cfg.ConfigDir(), "state"), data, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Could not write state")
		return
	}
}

func loadState() (State, ArmingMode, AlarmType) {
	data, err := os.ReadFile(path.Join(cfg.ConfigDir(), "state"))
	if err != nil {
		log.Error().Err(err).Msg("Could not read state")
		return "", "", ""
	}

	var s systemstate
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Error().Err(err).Msg("Could parse state")
		return "", "", ""
	}

	return s.State, s.Mode, s.Alarm
}
