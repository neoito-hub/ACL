package update_user

import "gorm.io/datatypes"

type RequestObject struct {
	UserID         string      `json:"user_id"`
	DeletedRoleIds []string    `json:"deleted_roles"`
	DeletedTeamIds []string    `json:"deleted_teams"`
	UserDetails    UserDetails `json:"user_details"`
}

type UserDetails struct {
	FullName string `json:"full_name"`
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
	SpaceID               string         `json:"space_id"`
	Name                  string         `json:"name"`
	Type                  string         `json:"type"`
	Email                 string         `json:"email"`
	Country               string         `json:"country"`
	BusinessName          string         `json:"business_name"`
	Address               string         `json:"address"`
	BusinessCategory      string         `json:"business_category"`
	Description           string         `json:"description"`
	MarketPlaceID         string         `json:"market_place_id"`
	DeveloperPortalAccess bool           `json:"developer_portal_access"`
	MetaData              datatypes.JSON `json:"meta_data"`
	LogoURL               string         `json:"logo_url"`
}
