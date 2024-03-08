package get_invite_by_id

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {
	// Validating request body and method
	// b, validateErr := ValidateRequest(w, r)
	// if validateErr != nil {
	// 	fmt.Printf("Error: %v\n", validateErr)
	// 	RespondWithError(w, http.StatusBadRequest, validateErr.Error())

	// 	return
	// }

	// // Validating and retreving user id from user access token
	// _, shieldVerifyError := VerifyAndGetUser(w, r)
	// if shieldVerifyError != nil {
	// 	fmt.Printf("shieldVerifyError: %v\n", shieldVerifyError)
	// 	RespondWithError(w, http.StatusUnauthorized, shieldVerifyError.Error())

	// 	return
	// }

	// db := DBInit()
	// sqlDB, dberr := db.DB()

	// if dberr != nil {
	// 	log.Fatalln(dberr)
	// }
	// defer sqlDB.Close()

	// fmt.Printf("b: %v\n", b)

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db

	shieldUser := ShieldUserData{
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}

	// var ResponseDetails ResponseData

	var inviteData InviteData
	var inviteDetails []InviteDetails
	var responseData ResponseData
	var validEmail Exists

	// // TODO

	res := db.Raw("SELECT i.invite_type, i.status, id.email, i.expires_at  FROM invites i INNER JOIN invite_details id ON id.invite_id = i.id where i.id=? group by i.invite_type, i.status, id.email, i.expires_at", b.InviteId).Scan(&inviteData)

	if res.Error != nil {
		// RespondWithError(w, http.StatusInternalServerError, "Server Error!")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	if res.RowsAffected < 1 {
		// RespondWithError(w, http.StatusNoContent, "NO RECORD FOUND!")
		// return
		handlerResp = common_services.BuildErrorResponse(true, "", Response{Err: true, Msg: "Invite details were fetched successfully.", Data: ResponseData{InviteType: inviteData.InviteType, Status: inviteData.Status, Msg: "Invite link you're trying to use is invalid or has been revoked."}}, http.StatusOK)
		return handlerResp
	}

	if inviteData.InviteType == 1 {

		if err := db.Raw(`select exists(select user_id from users where user_id = ? AND email=? limit 1)`, shieldUser.UserID, inviteData.Email).Scan(&validEmail).Error; err != nil {
			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)

			return handlerResp

		}

		// msg included im data part, http status changed to 200
		if !validEmail.Exists {
			handlerResp = common_services.BuildErrorResponse(true, "", Response{Err: true, Msg: "Invite details were fetched successfully.", Data: ResponseData{InviteType: inviteData.InviteType, Msg: "It looks like the invite link you're trying to use was created for another user."}}, http.StatusOK)

			return handlerResp
		}
	}

	if inviteData.Status == 2 {
		handlerResp = common_services.BuildErrorResponse(true, "", Response{Err: true, Msg: "Invite details were fetched successfully.", Data: ResponseData{InviteType: inviteData.InviteType, Status: inviteData.Status, Msg: "You have already accepted this invite."}}, http.StatusOK)
		return handlerResp
	}

	if inviteData.Status == 3 {
		handlerResp = common_services.BuildErrorResponse(true, "", Response{Err: true, Msg: "Invite details were fetched successfully.", Data: ResponseData{InviteType: inviteData.InviteType, Status: inviteData.Status, Msg: "This invite link has been declined."}}, http.StatusOK)
		return handlerResp
	}

	if time.Now().After(inviteData.ExpiresAt) {
		handlerResp = common_services.BuildErrorResponse(true, "", Response{Err: true, Msg: "Invite details were fetched successfully.", Data: ResponseData{InviteType: inviteData.InviteType, Status: inviteData.Status, Expired: true, Msg: "The invite link you were trying to use has expired."}}, http.StatusOK)

		return handlerResp
	}

	res = db.Raw("select s.space_id, s.name as space_name, json_agg(json_build_object('team_id',t.team_id,'team_name',t.name)) FILTER (WHERE t.team_id IS NOT NULL) as team_data, json_agg(json_build_object('role_id',r.id,'role_name',r.name)) FILTER (WHERE r.id IS NOT NULL) as role_data from invite_details id inner join spaces s on s.space_id = id.invited_space_id left join teams t on t.team_id = id.invited_team_id left join roles r on r.id = id.invited_role_id where id.invite_id = ? group by s.space_id", b.InviteId).Scan(&inviteDetails)

	if res.Error != nil {
		// RespondWithError(w, http.StatusInternalServerError, "Server Error!")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	if res.RowsAffected < 1 {
		// RespondWithError(w, http.StatusNoContent, "NO RECORD FOUND!")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "NO RECORD FOUND!", Response{}, http.StatusNoContent)
		return handlerResp
	}

	responseData.InviteDetails = inviteDetails
	responseData.InviteType = inviteData.InviteType
	responseData.Status = inviteData.Status
	responseData.Email = inviteData.Email
	responseData.ExpiresAt = inviteData.ExpiresAt
	responseData.Msg = fmt.Sprintf("You have been invited to collaborate in %s space", inviteDetails[0].SpaceName)

	var resp Response
	resp.Data = responseData
	resp.Err = false
	resp.Msg = "Invite details fetched successfully!"

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "", resp, http.StatusOK)
	return handlerResp
}
