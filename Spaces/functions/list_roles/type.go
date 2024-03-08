package list_roles

import "gorm.io/datatypes"

type RequestObject struct {
	SpaceID            string `json:"space_id"`
	SearchKeyword      string `json:"search_keyword"`
	PageLimit          int    `json:"page_limit"`
	Offset             int    `json:"offset"`
	DisplayMemberCount int    `json:"display_member_count"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResultData struct {
	TotalCount int             `json:"total_count"`
	Data       []RolesListData `json:"data"`
}

type RolesListData struct {
	RoleID      string         `json:"role_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	SpaceID     string         `json:"space_id"`
	MemberCount int            `json:"member_count"`
	Members     datatypes.JSON `json:"members"`
	IsOwner     bool           `json:"is_owner"`
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
