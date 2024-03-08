package teams_search_user

type RequestObject struct {
	SearchString string `json:"search_string"`
	TeamID       string `json:"team_id"`
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
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}
