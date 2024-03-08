package create_team

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aidarkhanov/nanoid"
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
	// userData, shieldVerifyError := VerifyAndGetUser(w, r)
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

	userData := ShieldUserData{
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}

	tx := db.Begin()

	// team and member common ID
	TeamID := nanoid.New()

	// Create member
	member := models.Member{
		ID:   TeamID,
		Type: "T",
	}

	memberResult := tx.Create(&member)

	if memberResult.Error != nil {
		log.Println("!!! Error creating member => ", memberResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating member")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating member", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	// Create Team
	team := models.Team{
		TeamID:      TeamID,
		OwnerID:     b.SpaceID,
		UpdatedBy:   userData.UserID,
		Name:        b.Name,
		Description: b.Description,
	}

	teamResult := tx.Create(&team)

	if teamResult.Error != nil {
		log.Println("!!! Error creating team => ", teamResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating team")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating team", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	// inserting to member taable as team owner

	teamMember := models.TeamMember{
		ID:          nanoid.New(),
		OwnerTeamID: TeamID,
		MemberID:    userData.UserID,
		IsOwner:     true,
	}

	teamMemberResult := tx.Create(&teamMember)

	if teamMemberResult.Error != nil {
		log.Println("!!! Error creating team member=> ", teamMemberResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating team member")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating team member", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{Err: false, Msg: "Successfully created team"})

	handlerResp = common_services.BuildResponse(false, "Successfully created team", Response{Err: false, Msg: "Successfully created team"}, http.StatusOK)
	return handlerResp
}
