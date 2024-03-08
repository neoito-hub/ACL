package search_user

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

	var userDetails []UserData

	searchParam := "%" + b.SearchString + "%"
	// TODO
	db.Raw(`SELECT u.user_id,u.user_name,u.email FROM users u left join space_members sm on (sm.owner_user_id = u.user_id 
		and sm.owner_space_id=?) left join invites i on i.email=u.email 
		where sm.id is null and i.id is null and (u.user_name ilike ? or u.email ilike ?)`, b.SpaceID, searchParam, searchParam).Scan(&userDetails)

	var resp Response
	resp.Data = userDetails
	resp.Err = false
	resp.Msg = "User details retrieved successfully!"

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "User details retrieved successfully!", resp, http.StatusOK)
	return handlerResp
}
