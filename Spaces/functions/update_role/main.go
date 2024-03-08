package update_role

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

	// Update Role data
	var roleDetail models.Role

	res := db.Model(&models.Role{}).Where("id=?", b.RoleID).Updates(map[string]interface{}{"name": b.Name, "description": b.Description, "updated_by": userData.UserID}).Scan(&roleDetail)

	if res.Error != nil {
		// RespondWithError(w, http.StatusBadRequest, "Error on updating role")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error on updating role", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	if res.RowsAffected < 1 {
		// RespondWithError(w, http.StatusNoContent, "NO RECORD FOUND!")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "NO RECORD FOUND!", Response{}, http.StatusNoContent)
		return handlerResp
	}

	// RespondWithJSON(w, http.StatusOK, Response{
	// 	Data: roleDetail,
	// 	Err:  false,
	// 	Msg:  "role detail updated successfully!",
	// })

	handlerResp = common_services.BuildResponse(false, "role detail updated successfully!", Response{
		Data: roleDetail,
		Err:  false,
		Msg:  "role detail updated successfully!",
	}, http.StatusOK)
	return handlerResp
}
