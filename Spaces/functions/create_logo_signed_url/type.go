package create_logo_signed_url

type ReturnObject struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type RequestObject struct {
	FileExtension string `json:"file_extension"`
}

type Response struct {
	Err  bool        `json:"err"`  // m
	Msg  string      `json:"msg"`  // m
	Data interface{} `json:"data"` // m
}

type ErrResponse struct {
	Err bool   `json:"err"` // m
	Msg string `json:"msg"` // m
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
