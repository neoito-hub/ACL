package get_space_by_id

import "gorm.io/datatypes"

type RequestObject struct {
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

type SpaceData struct {
	SpaceID               string         `json:"space_id"`
	LegalID               string         `json:"legal_id"`
	Type                  string         `json:"type"`
	Name                  string         `json:"name"`
	BusinessName          string         `json:"business_name"`
	Address               string         `json:"address"`
	Email                 string         `json:"email"`
	Country               string         `json:"country"`
	BusinessCategory      string         `json:"business_category"`
	Description           string         `json:"description"`
	MetaData              datatypes.JSON `json:"metadata"`
	TaxPayerID            string         `json:"tax_payer_id"`
	DistinguishedName     string         `json:"distinguished_name"`
	Status                int            `json:"status"`
	MarketPlaceID         string         `json:"market_place_id"`
	DeveloperPortalAccess bool           `json:"developer_portal_access"`
}
