package create_entity

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	tx := db.Begin()

	// team and member common ID
	EntityID := nanoid.New()

	// Create member
	entity := models.Entities{
		EntityID:  EntityID,
		Type:      int64(b.TypeID),
		Label:     b.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entityResult := tx.Create(&entity)

	if entityResult.Error != nil {
		log.Println("!!! Error creating member => ", entityResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating member")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating entity", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	// Create Team
	entitySpaceMapping := models.EntitySpaceMapping{
		ID:            nanoid.New(),
		OwnerEntityID: EntityID,
		OwnerSpaceID:  b.SpaceID,
	}

	entitySpaceMappingResult := tx.Create(&entitySpaceMapping)

	if entitySpaceMappingResult.Error != nil {
		log.Println("!!! Error creating entity space mapping => ", entitySpaceMappingResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating team")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating Entity", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{Err: false, Msg: "Successfully created team"})

	handlerResp = common_services.BuildResponse(false, "Successfully created Entity", Response{Err: false, Msg: "Successfully created Entity"}, http.StatusOK)
	return handlerResp
}
