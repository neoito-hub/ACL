package accept_invite

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/neoito-hub/ACL-Block/Data-Models/models"
	"github.com/neoito-hub/ACL-Block/spaces/common_services"
	"gorm.io/gorm"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {
	// Validating request body and method

	// b, validateErr := ValidateRequest(w, r)
	// if validateErr != nil {
	// 	fmt.Printf("Error: %v\n", validateErr)
	// 	RespondWithError(w, http.StatusBadRequest, validateErr.Error())

	// 	return
	// }

	// Validating and retreving user id from user access token
	// shieldUser, shieldVerifyError := VerifyAndGetUser(w, r)
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

	var tx = db.Begin()

	var inviteData InviteData
	var validEmail Exists

	if err := tx.Raw(`select status, expires_at, invite_type, email from invites where id = ?`, b.InviteID).Scan(&inviteData).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)

		return handlerResp
	}

	if inviteData.InviteType == 1 {

		if err := tx.Raw(`select exists(select user_id from users where user_id = ? AND email=? limit 1)`, shieldUser.UserID, inviteData.Email).Scan(&validEmail).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return
			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)

			return handlerResp

		}

		if !validEmail.Exists {
			tx.Rollback()
			// RespondWithError(w, http.StatusUnauthorized, "Unauthorized Access")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Unauthorized Access", Response{Err: true, Msg: "Unauthorized Access"}, http.StatusUnauthorized)

			return handlerResp
		}
	}

	if time.Now().After(inviteData.ExpiresAt) {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Link Expired")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Link Expired", Response{Err: true, Msg: "Link Expired"}, http.StatusBadRequest)

		return handlerResp
	}

	if inviteData.Status == 2 || inviteData.Status == 3 {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Completed or Declined Link")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Completed or Declined Link", Response{Err: true, Msg: "Completed or Declined Link"}, http.StatusBadRequest)

		return handlerResp
	}

	err, createResponse := AddUserToSpaceTeam(tx, &shieldUser, b, inviteData.InviteType)

	if err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)

		return handlerResp
	}

	tx.Commit()
	var resp Response
	resp.Data = createResponse
	resp.Err = false
	resp.Msg = "Space details updated successfully!"

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Space details updated successfully!", resp, http.StatusOK)

	return handlerResp
}

func AddUserToSpaceTeam(tx *gorm.DB, shieldUser *ShieldUserData, inviteData RequestObject, inviteType int) (error, ResponseData) {

	var spaceResponse []UserInsertSpaceData
	var teamResponse []UserInsertTeamData
	var roleResponse []UserInsertRoleData

	if inviteType == 1 { // email
		if err := tx.Raw(`insert into space_members(owner_user_id,owner_space_id,opt_counter,id,created_at,updated_at) SELECT u.user_id, id.invited_space_id, 0, nanoid(),now(),now() FROM invite_details id INNER JOIN users u ON u.email = id.email WHERE id.invite_id = ? AND id.invited_space_id NOT IN(SELECT sm.owner_space_id FROM space_members sm INNER JOIN invite_details id ON id.invited_space_id = sm.owner_space_id WHERE sm.owner_user_id = ?) group by u.user_id, id.invited_space_id returning *`, inviteData.InviteID, shieldUser.UserID).Scan(&spaceResponse).Error; err != nil {
			return err, ResponseData{}
		}

		if err := tx.Raw(`insert into team_members(member_id,owner_team_id,opt_counter,id,created_at,updated_at) SELECT u.user_id, id.invited_team_id, 0, nanoid(),now(),now() FROM invite_details id INNER JOIN users u ON u.email = id.email WHERE id.invite_id = ? AND coalesce(invited_team_id, '') !=  '' AND id.invited_team_id NOT IN(SELECT tm.owner_team_id FROM team_members tm INNER JOIN invite_details id ON id.invited_team_id = tm.owner_team_id WHERE tm.member_id = ?) group by u.user_id, id.invited_team_id  returning *`, inviteData.InviteID, shieldUser.UserID).Scan(&teamResponse).Error; err != nil {
			return err, ResponseData{}
		}

		if err := tx.Raw(`insert into member_roles(owner_user_id,owner_space_id,role_id,opt_counter,id,created_at,updated_at) SELECT u.user_id, id.invited_space_id, id.invited_role_id, 0, nanoid(),now(),now() FROM invite_details id INNER JOIN users u ON u.email = id.email WHERE id.invite_id = ? AND coalesce(invited_role_id, '') !=  '' AND id.invited_role_id NOT IN(SELECT mr.role_id FROM member_roles mr INNER JOIN invite_details id ON id.invited_role_id = mr.role_id WHERE mr.owner_user_id = ?) group by u.user_id, id.invited_space_id, id.invited_role_id  returning *`, inviteData.InviteID, shieldUser.UserID).Scan(&roleResponse).Error; err != nil {
			return err, ResponseData{}
		}

		// update invite status

		if err := tx.Model(&models.Invites{}).Where("id = ?", inviteData.InviteID).Update("status", 2).Error; err != nil {
			return err, ResponseData{}
		}

	} else { // link
		if err := tx.Raw(`insert into space_members(owner_user_id,owner_space_id,opt_counter,id,created_at,updated_at) SELECT ?, id.invited_space_id, 0, nanoid(),now(),now() FROM invite_details id WHERE id.invite_id = ? AND id.invited_space_id NOT IN(SELECT sm.owner_space_id FROM space_members sm INNER JOIN invite_details id ON id.invited_space_id = sm.owner_space_id WHERE sm.owner_user_id = ?) group by id.invited_space_id returning *`, shieldUser.UserID, inviteData.InviteID, shieldUser.UserID).Scan(&spaceResponse).Error; err != nil {
			return err, ResponseData{}
		}

		if err := tx.Raw(`insert into team_members(member_id,owner_team_id,opt_counter,id,created_at,updated_at) SELECT ?, id.invited_team_id, 0, nanoid(),now(),now() FROM invite_details id WHERE id.invite_id = ? AND coalesce(invited_team_id, '') !=  '' AND id.invited_team_id NOT IN(SELECT tm.owner_team_id FROM team_members tm INNER JOIN invite_details id ON id.invited_team_id = tm.owner_team_id WHERE tm.member_id = ?) group by id.invited_team_id returning *`, shieldUser.UserID, inviteData.InviteID, shieldUser.UserID).Scan(&teamResponse).Error; err != nil {
			return err, ResponseData{}
		}

		if err := tx.Raw(`insert into member_roles(owner_user_id,owner_space_id,role_id,opt_counter,id,created_at,updated_at) SELECT ?, id.invited_space_id, id.invited_role_id, 0, nanoid(),now(),now() FROM invite_details id WHERE id.invite_id = ? AND coalesce(invited_role_id, '') !=  '' AND id.invited_role_id NOT IN(SELECT mr.role_id FROM member_roles mr INNER JOIN invite_details id ON id.invited_role_id = mr.role_id WHERE mr.owner_user_id = ?) group by id.invited_space_id, id.invited_role_id returning *`, shieldUser.UserID, inviteData.InviteID, shieldUser.UserID).Scan(&roleResponse).Error; err != nil {
			return err, ResponseData{}
		}
	}

	return nil, ResponseData{
		UserInsertSpaceData: spaceResponse,
		UserInsertTeamData:  teamResponse,
		UserInsertRoleData:  roleResponse,
	}

}
