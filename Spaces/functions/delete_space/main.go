package delete_space

import (
	"encoding/json"
	"net/http"

	"github.com/neoito-hub/ACL-Block/Data-Models/models"
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

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db

	tx := db.Begin()

	var acPolGrpIds []models.AcPolGrp
	res := tx.Raw("select id from ac_pol_grps where member_id = ?", b.ID).Scan(&acPolGrpIds)

	if res.Error != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Invalid request payload")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	var teamIds []*models.Team
	teamRes := tx.Raw("select team_id from teams where owner_id = ?", b.ID).Scan(&teamIds)

	if teamRes.Error != nil {
		tx.Rollback()

		// RespondWithError(w, http.StatusBadRequest, "Invalid request payload")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	teamIdsList := make([]string, len(teamIds))

	for i, t := range teamIds {
		teamIdsList[i] = t.TeamID
	}

	acPolGrpIdsList := make([]string, len(acPolGrpIds))

	for i, t := range acPolGrpIds {
		acPolGrpIdsList[i] = t.ID
	}

	acPolGrpsErr := DeleteQueryRun(tx, "DELETE FROM ac_pol_grps WHERE member_id = ?", b.ID)
	if acPolGrpsErr != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error deleting ac_pol_grps ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting ac_pol_grps", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	polGpPoliciesErr := DeleteQueryRun(tx, "DELETE FROM pol_gp_policies WHERE ac_pol_grp_id in (?)", acPolGrpIdsList)
	if polGpPoliciesErr != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error deleting pol_gp_policies")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting pol_gp_policies", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	rolesErr := DeleteQueryRun(tx, "DELETE FROM roles WHERE owner_space_id = ?", b.ID)
	if rolesErr != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error deleting roles ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting roles", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	memberRolesErr := DeleteQueryRun(tx, "DELETE FROM member_roles WHERE owner_space_id = ?", b.ID)
	if memberRolesErr != nil {
		tx.Rollback()

		// RespondWithError(w, http.StatusBadRequest, "Error deleting member_roles ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting member_roles", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	teamsErr := DeleteQueryRun(tx, "DELETE FROM teams WHERE owner_id = ?", b.ID)
	if teamsErr != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error deleting teams ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting teams", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	teamMembersErr := DeleteQueryRun(tx, "DELETE FROM team_members WHERE space_id in (?)", teamIdsList)
	if teamMembersErr != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error deleting team_members ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting team_members", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	spacesErr := DeleteQueryRun(tx, "DELETE FROM spaces WHERE space_id = ?", b.ID)
	if spacesErr != nil {
		tx.Rollback()

		// RespondWithError(w, http.StatusBadRequest, "Error deleting spaces ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting spaces", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	membersErr := DeleteQueryRun(tx, "DELETE FROM members WHERE id = ?", b.ID)
	if membersErr != nil {
		tx.Rollback()

		// RespondWithError(w, http.StatusBadRequest, "Error deleting space member ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting space member", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{
	// 	Err: false,
	// 	Msg: "Space deleted successfully!",
	// })

	handlerResp = common_services.BuildResponse(false, "Space deleted successfully!", Response{
		Err: false,
		Msg: "Space deleted successfully!",
	}, http.StatusOK)
	return handlerResp

}
