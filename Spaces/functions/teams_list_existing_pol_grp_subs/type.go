package teams_list_existing_pol_grp_subs

import "time"

type RequestObject struct {
	SpaceID       string `json:"space_id"`
	TeamID        string `json:"team_id"`
	SearchKeyword string `json:"search_keyword"`
	PageLimit     int    `json:"page_limit"`
	Offset        int    `json:"offset"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResultData struct {
	TotalCount int                `json:"total_count"`
	Data       []PoliciesListData `json:"data"`
}

type PoliciesListData struct {
	SubsID         string    `json:"subs_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	IsPredefined   bool      `json:"is_predefined"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PermissionName string    `json:"permission_name"`
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
