package user_list_entities

import "gorm.io/datatypes"

type RequestObject struct {
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
	SortColumn    string     `json:"sort_column"`
	SortDirection string     `json:"sort_direction"`
	IsFiltersSet  bool       `json:"is_filters_set"` // m
	Conditions    Conditions `json:"conditions"`
	UserID        string     `json:"user_id"`
	TypeID        int        `json:"type_id"`
}

type Conditions struct {
	IsKeywordSearch bool   `json:"is_keyword_search"`
	Keyword         string `json:"keyword"`
}

type EntityData struct {
	Entities []Entity `json:"entities"`
	Count    int      `json:"count"`
}

type Entity struct {
	EntityID     string         `json:"entity_id"`
	Label        string         `json:"label"`
	Type         int            `json:"type"` // app - 1 , ui-container - 2,ui-elements 3  fn - 4, data - 5,function shared block -6
	PolicyGroups datatypes.JSON `json:"policy_groups"`
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

type TagsCreateResponse struct {
	ID      string `json:"id"`
	TagName string `json:"tag_name"`
}

type BlockTagsMappingCreateResponse struct {
	ID      string `json:"id"`
	TagID   string `json:"tag_id"`
	BlockID string `json:"block_id"`
}
