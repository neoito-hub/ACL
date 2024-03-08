package create_space

type RequestObject struct {
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	Email                 string `json:"email"`
	Country               string `json:"country"`
	BusinessName          string `json:"business_name"`
	Address               string `json:"address"`
	BusinessCategory      string `json:"business_category"`
	Description           string `json:"description"`
	MarketPlaceID         string `json:"market_place_id"`
	LogoURL               string `json:"logo_url"`
	DeveloperPortalAccess bool   `json:"developer_portal_access"`
}

type subExists struct {
	PolGrpSubsExists bool `json:"pol_grp_subs_exists"`
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
