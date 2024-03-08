package check_assigned_block_to_app

type RequestObject struct {
	AppID   string `json:"app_id"`
	BlockID string `json:"block_id"`
	SpaceID string `json:"space_id"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ExistResponse struct {
	Exist         bool `json:"exist"`
	InUse         bool `json:"in_use"`
	CanAssign     bool `json:"can_assign"`
	InUseCount    int  `json:"in_use_count"`
	PurchaseCount int  `json:"purchase_count"`
	// CanReAssign bool `json:"can_re_assign"`
}
