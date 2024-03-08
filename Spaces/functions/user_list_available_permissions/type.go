package user_list_available_permissions

import (
	"gorm.io/datatypes"
)

type RequestObject struct {
	UserID        string `json:"user_id"`
	SearchKeyword string `json:"search_keyword"`
	PageLimit     int    `json:"page_limit"`
	Offset        int    `json:"offset"`
	SortColumn    string `json:"sort_column"`
	SortDirection string `json:"sort_direction"`
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
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	PermissionID string         `json:"permission_id"`
	PolicyGroups datatypes.JSON `json:"policy_groups"`
	PgCount      int            `json:"pg_count"`
	EntityTypes  datatypes.JSON `json:"entity_types"`
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
