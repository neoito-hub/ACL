package resend_invite_email

type RequestObject struct {
	InviteIds []string `json:"invite_ids"`
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

type ExistingInvites struct {
	InviteID    string `json:"invite_id"`
	Email       string `json:"email"`
	NewInviteID string `json:"new_invite_id"`
	SpaceID     string `json:"space_id"`
	SpaceName   string `json:"space_name"`
}

type LinkPayload struct {
	InviteID string `json:"invite_id"`
}

type EmailSendPayload struct {
	Link  string `json:"link"`
	Email string `json:"email"`
}

type EmailStruct struct {
	InviteLink  string
	SpaceName   string
	Email       string
	FirstLetter string
}
