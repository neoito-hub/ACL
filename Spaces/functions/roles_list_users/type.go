package roles_list_users

import "gorm.io/datatypes"

type RequestObject struct {
	RoleID          string     `json:"role_id"`
	IsFilter        bool       `json:"is_filter"`         // m
	IsKeywordSearch bool       `json:"is_keyword_search"` // m
	Conditions      Conditions `json:"conditions"`        // m
	PageLimit       int64      `json:"page_limit"`        // m
	Offset          int        `json:"offset"`            // m
	Active          string     `json:"active"`            // m
	Direction       string     `json:"direction"`         // m

}

type Conditions struct {
	Keyword string `json:"search_keyword"` // o
	Filter  Filter `json:"filter"`         // o
}

type Filter struct {
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

type ResultData struct {
	TotalCount int    `json:"total_count"`
	Data       []Data `json:"data"`
}
type Data struct {
	MemberRoleID string         `json:"member_role_id"`
	UserID       string         `json:"user_id"`
	UserName     string         `json:"user_name"`
	FullName     string         `json:"full_name"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	UpdatedDate  string         `json:"updated_date"`
	CreatedDate  string         `json:"created_date"`
	Type         int            `json:"type"`
	Teams        datatypes.JSON `json:"teams"`
}
