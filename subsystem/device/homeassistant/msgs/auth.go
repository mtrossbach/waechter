package msgs

type AuthResponse struct {
	Type      MsgType `json:"type"`
	HaVersion string  `json:"ha_version"`
}

type AuthRequest struct {
	Type        MsgType `json:"type"`
	AccessToken string  `json:"access_token"`
}
