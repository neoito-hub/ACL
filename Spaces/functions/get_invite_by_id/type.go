package get_invite_by_id

import (
	"time"

	"gorm.io/datatypes"
)

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

type ResponseData struct {
	InviteDetails []InviteDetails `json:"invite_details"`
	InviteType    int             `json:"invite_type"` //1-email invite 2 invite link
	Status        int             `json:"status"`
	Email         string          `json:"email"`
	ExpiresAt     time.Time       `json:"expires_at"`
	Expired       bool            `json:"expired"`
	Msg           string          `json:"msg"`
}

type InviteDetails struct {
	SpaceID   string         `json:"space_id"`
	SpaceName string         `json:"space_name"`
	TeamData  datatypes.JSON `json:"team_data"`
	RoleData  datatypes.JSON `json:"role_data"`
}

type TeamData struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
}

type RoleData struct {
	RoleID   string `json:"role_id"`
	RoleName string `json:"role_name"`
}

type InviteData struct {
	InviteType int       `json:"invite_type"`
	Status     int       `json:"status"`
	Email      string    `json:"email"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type Exists struct {
	Exists bool `json:"exists"`
}
