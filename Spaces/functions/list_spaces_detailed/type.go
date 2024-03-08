package list_spaces_detailed

import "gorm.io/datatypes"

type RequestObject struct {
	State         int32  `json:"state"` // 0:all, 1:personal, 2:business
	SearchKeyword string `json:"search_keyword"`
	PageLimit     int    `json:"page_limit"`
	Offset        int    `json:"offset"`
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

type SpaceDetails struct {
	SpaceID     string         `json:"space_id"`
	SpaceName   string         `json:"space_name"`
	LogoURL     string         `json:"logo_url"`
	Type        string         `json:"type"`
	CreatedAt   datatypes.Date `json:"created_at"`
	MemberCount int            `json:"member_count"`
	EntityCount    int         `json:"entity_count"`
	MySpace     bool           `json:"my_space"`
	OwnerUserID string         `json:"owner_user_id"`
	FullName    string         `json:"full_name"`
	Email       string         `json:"email"`
	UserName    string         `json:"user_name"`
}

type ResultData struct {
	TotalCount int            `json:"total_count"`
	Data       []SpaceDetails `json:"data"`
}
