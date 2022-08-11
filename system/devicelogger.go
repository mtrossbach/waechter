package system

import "github.com/rs/zerolog"

func DevLog(device Device, e *zerolog.Event) *zerolog.Event {
	return e.Str("id", device.GetId()).Str("displayName", device.GetDisplayName()).Str("type", string(device.GetType())).Str("device", device.GetSubsystem())
}
