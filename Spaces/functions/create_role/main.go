package create_role

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

	roleID := nanoid.New()
	// Create Role
	role := models.Role{
		ID:           roleID,
		Name:         b.Name,
		Description:  b.Description,
		OwnerSpaceID: b.SpaceID,
		CreatedBy:    userData.UserID,
		UpdatedBy:    userData.UserID,
	}

	roleResult := tx.Create(&role)

	if roleResult.Error != nil {
		log.Println("!!! Error creating role => ", roleResult.Error)
		tx.Rollback()

		// RespondWithError(w, http.StatusBadRequest, "Error creating role")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating role", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	memberRole := models.MemberRole{
		ID:           nanoid.New(),
		OwnerUserID:  userData.UserID,
		RoleID:       roleID,
		OwnerSpaceID: b.SpaceID,
	}

	memberRoleResult := tx.Create(&memberRole)

	if memberRoleResult.Error != nil {
		log.Println("!!! Error creating member role=> ", memberRoleResult.Error)
		tx.Rollback()

		// RespondWithError(w, http.StatusBadRequest, "Error creating member role")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating member role", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{Err: false, Msg: "Role added successfully"})

	handlerResp = common_services.BuildResponse(false, "Role added successfully", Response{Err: false, Msg: "Role added successfully"}, http.StatusOK)
	return handlerResp
}
