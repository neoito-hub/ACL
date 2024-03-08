package user_list_available_entities

import (
	"time"
)

type RequestObject struct {
	SearchKeyword string `json:"search_keyword"`
	PageLimit     int    `json:"page_limit"`
	Offset        int    `json:"offset"`
	SortColumn    string `json:"sort_column"`
	SortDirection string `json:"sort_direction"`
	EntityTypes   []int  `json:"entity_types"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResultData struct {
	TotalCount int                `json:"total_count"`
	Data       []EntitiesListData `json:"data"`
}

type EntitiesListData struct {
	EntityID  string    `json:"entity_id"`
	Type      int       `json:"type"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
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
