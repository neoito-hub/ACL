package delete_team

import (
	"encoding/json"
	"net/http"

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

	teamsErr := DeleteQueryRun(tx, "DELETE FROM teams WHERE team_id=?", b.TeamID)
	if teamsErr != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error deleting teams ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting teams", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	teamMembersErr := DeleteQueryRun(tx, "DELETE FROM team_members WHERE owner_team_id=?", b.TeamID)
	if teamMembersErr != nil {
		tx.Rollback()

		// RespondWithError(w, http.StatusBadRequest, "Error deleting team_members ")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error deleting team_members", Response{}, http.StatusBadRequest)
		return handlerResp
	}

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
