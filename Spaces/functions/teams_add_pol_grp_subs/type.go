package teams_add_pol_grp_subs

type RequestObject struct {
	SpaceID     string   `json:"space_id"`
	TeamID      string   `json:"team_id"`
	AcPolGrpIDs []string `json:"ac_pol_grp_ids"`
}

type ResponseData struct {
	ID         string `json:"id"`
	SpaceID    string `json:"space_id"`
	TeamID     string `json:"team_id"`
	AcPolGrpID string `json:"ac_pol_grp_id"`
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
