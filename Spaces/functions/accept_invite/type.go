package accept_invite

import "time"

type RequestObject struct {
	InviteID string `json:"invite_id"`
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

type UserInsertSpaceData struct {
	OwnerSpaceID string `json:"owner_space_id"`
	OwnerUserID  string `json:"owner_user_id"`
}

type UserInsertTeamData struct {
	OwnerTeamID string `json:"owner_team_id"`
	MemberID    string `json:"member_id"`
}

type UserInsertRoleData struct {
	RoleID      string `json:"role_id"`
	OwnerUserID string `json:"owner_user_id"`
}

type ResponseData struct {
	UserInsertSpaceData []UserInsertSpaceData `json:"space_data"`
	UserInsertTeamData  []UserInsertTeamData  `json:"team_data"`
	UserInsertRoleData  []UserInsertRoleData  `json:"role_data"`
}

type InviteData struct {
	Status     int
	ExpiresAt  time.Time
	InviteType int
	Email      string
}

type Exists struct {
	Exists bool `json:"exists"`
}
