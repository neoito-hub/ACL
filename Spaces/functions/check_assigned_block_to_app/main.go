package check_assigned_block_to_app

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

	spaceID := b.SpaceID
	if len(spaceID) < 1 {
		spaceID = payload.SpaceID
	}

	// TODO: Block paid has paschased check
	purchasedRes := db.Exec(`select id from order_items where variant_block_id = ? and space_id = ?`, b.BlockID, spaceID)

	if purchasedRes.Error != nil {
		handlerResp = common_services.BuildErrorResponse(true, purchasedRes.Error.Error(), Response{
			Msg: purchasedRes.Error.Error(), Data: purchasedRes.Error, Err: true,
		}, http.StatusBadRequest)

		return handlerResp
	}

	purchasedCount := purchasedRes.RowsAffected

	if purchasedCount < 1 {
		handlerResp = common_services.BuildErrorResponse(true, "No purchase found", Response{
			Err: true,
			Msg: "No purchase found",
		}, http.StatusForbidden)

		return handlerResp
	}

	var existingAssigned []struct {
		ID     string `json:"id"`
		Status int    `json:"status"`
		AppID  string `json:"app_id"`
	}

	existingAssignedRes := db.Raw(`select id, status, app_id from block_app_assigns where block_id = ? and space_id = ?`, b.BlockID, spaceID).Scan(&existingAssigned)

	if existingAssignedRes.Error != nil {
		handlerResp = common_services.BuildErrorResponse(true, existingAssignedRes.Error.Error(), Response{
			Msg: existingAssignedRes.Error.Error(), Data: existingAssignedRes.Error, Err: true,
		}, http.StatusBadRequest)

		return handlerResp
	}

	var repData = ExistResponse{
		InUseCount:    0,
		InUse:         false,
		Exist:         false,
		PurchaseCount: int(purchasedCount),
		CanAssign:     false,
	}

	statusAssigned := 1
	statusInUse := 2

	for _, eA := range existingAssigned {
		if eA.Status == statusInUse {
			if eA.AppID == b.AppID {
				repData.InUse = true
				repData.Exist = true
			}
			repData.InUseCount++
		} else if eA.Status == statusAssigned {
			if eA.AppID == b.AppID {
				repData.Exist = true
			}
			// repData.CanReAssign = true
		}
	}

	if int64(repData.InUseCount) < purchasedCount {
		repData.CanAssign = true
	}

	var resp Response
	resp.Data = repData
	resp.Err = false
	resp.Msg = "Data retrieved!"

	handlerResp = common_services.BuildResponse(false, "Data retrieved!", resp, http.StatusOK)

	return handlerResp
}
