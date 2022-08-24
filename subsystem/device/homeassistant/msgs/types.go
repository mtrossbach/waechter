package msgs

type MsgType string

const (
	AuthRequired         MsgType = "auth_required"
	Auth                 MsgType = "auth"
	AuthOk               MsgType = "auth_ok"
	AuthInvalid          MsgType = "auth_invalid"
	Result               MsgType = "result"
	SubscribeEvents      MsgType = "subscribe_events"
	Event                MsgType = "event"
	SubscribeTrigger     MsgType = "subscribe_trigger"
	UnsubscribeEvents    MsgType = "unsubscribe_events"
	FireEvent            MsgType = "fire_event"
	CallService          MsgType = "call_service"
	GetStates            MsgType = "get_states"
	GetConfig            MsgType = "get_config"
	GetServices          MsgType = "get_services"
	GetPanels            MsgType = "get_panels"
	MediaPlayerThumbnail MsgType = "media_player_thumbnail"
	Ping                 MsgType = "ping"
	Pong                 MsgType = "pong"
	ValidateConfig       MsgType = "validate_config"
)
