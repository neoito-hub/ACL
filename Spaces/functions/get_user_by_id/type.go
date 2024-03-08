package get_user_by_id

import "gorm.io/datatypes"

type RequestObject struct {
	UserID  string `json:"user_id"`
	SpaceID string `json:"space_id"`
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

type UserData struct {
	Roles       datatypes.JSON `json:"roles"`
	Teams       datatypes.JSON `json:"teams"`
	UserID      string         `json:"user_id"`
	UserName    string         `json:"user_name"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	UpdatedDate string         `json:"updated_date"`
	CreatedDate string         `json:"created_date"`
	FullName    string         `json:"full_name"`
}
