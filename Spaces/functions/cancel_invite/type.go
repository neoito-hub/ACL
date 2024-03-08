package cancel_invite

type RequestObject struct {
	InviteId string `json:"invite_id"`
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

type UserInsertData struct {
	OwnerSpaceID string `json:"owner_space_id"`
	OwnerUserID  string `json:"owner_user_id"`
}

type MemberRoleData struct {
	SpaceID     string
	OwnerUserID string
	OptCOunter  int
	ID          string
}
