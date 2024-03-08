package create_invite_link

type RequestObject struct {
	Data []TeamPayload `json:"data"`
}

type TeamPayload struct {
	SpaceID string   `json:"space_id"`
	TeamIDs []string `json:"team_ids"`
}

type SpaceIds struct {
	SpaceID string `json:"space_id"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Err bool   `json:"err"`
	Msg string `json:"msg"`
}

type ShiledResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

type ShieldUserData struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type LinkResponse struct {
	InviteLink string `json:"invite_link"`
}

type LinkPayload struct {
	InviteID string `json:"invite_id"`
}
