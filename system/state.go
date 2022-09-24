package system

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"os"
	"path"
	"time"
)

type State struct {
	ArmState      ArmState   `json:"armState"`
	Alarm         AlarmType  `json:"alarm"`
	WrongPinCount int        `json:"wrongPinCount"`
	DelayEnd      *time.Time `json:"delayEnd"`
}

func (s *State) IsArmed() bool {
	return s.ArmState == ArmedStayState || s.ArmState == ArmedAwayState
}

func (s *State) IsArmedOrExitDelay() bool {
	return s.IsArmed() || s.ArmState == ExitDelayState
}

func (s *State) writeToDisk() {
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

func (s *State) loadFromDisk() {
	data, err := os.ReadFile(path.Join(cfg.ConfigDir(), "state"))
	if err != nil {
		log.Error().Err(err).Msg("Could not read state")
		return
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		log.Error().Err(err).Msg("Could parse state")
	}
}
