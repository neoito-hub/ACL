package unassign_block_from_app

import (
	"encoding/json"
	"net/http"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {
	var b RequestObject

	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	db := payload.Db

	var isInUse struct {
		Status int `json:"status"`
	}

	isInUseRes := db.Raw(`select status from block_app_assigns where id = ?`, b.BlockAppAssignID).Scan(&isInUse)

	if isInUseRes.Error != nil {
		handlerResp = common_services.BuildErrorResponse(true, isInUseRes.Error.Error(), Response{
			Msg: isInUseRes.Error.Error(), Data: isInUseRes.Error, Err: true,
		}, http.StatusBadRequest)

		return handlerResp
	}

	if isInUse.Status == 2 {
		handlerResp = common_services.BuildErrorResponse(true, "Block is in use. No changes are allowed", Response{
			Msg: "Block is in use. No changes are allowed", Err: true,
		}, http.StatusBadRequest)

		return handlerResp
	}

	createRes := db.Exec(`delete from block_app_assigns where id = ?`, b.BlockAppAssignID)

	if createRes.Error != nil {
		handlerResp = common_services.BuildErrorResponse(true, createRes.Error.Error(), Response{
			Msg: createRes.Error.Error(), Data: createRes.Error, Err: true,
		}, http.StatusBadRequest)

		return handlerResp
	}

	var resp Response
	resp.Data = nil
	resp.Err = false
	resp.Msg = "Block unassinged successfully!"

	handlerResp = common_services.BuildResponse(false, "Block unassinged successfully!", resp, http.StatusOK)

	return handlerResp
}
