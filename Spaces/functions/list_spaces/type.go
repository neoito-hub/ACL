package list_spaces

import "gorm.io/datatypes"

type RequestObject struct {
	SearchKeyword string `json:"search_keyword"`
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
	SpaceID   string         `json:"space_id"`
	SpaceName string         `json:"space_name"`
	LogoURL   string         `json:"logo_url"`
	Type      string         `json:"type"`
	IsDefault bool           `json:"is_default"`
	Roles     datatypes.JSON `json:"roles"`
}
