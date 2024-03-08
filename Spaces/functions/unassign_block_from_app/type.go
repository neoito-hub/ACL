package unassign_block_from_app

type RequestObject struct {
	BlockAppAssignID string `json:"block_app_assign_id"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
