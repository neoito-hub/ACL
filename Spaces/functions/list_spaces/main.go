package list_spaces

import (
	"encoding/json"
	"net/http"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {
	// Validating request body and method
	// validateErr := ValidateRequest(w, r)
	// if validateErr != nil {
	// 	fmt.Printf("Error: %v\n", validateErr)
	// 	RespondWithError(w, http.StatusBadRequest, validateErr.Error())

	// 	return
	// }

	// // Validating and retreving user id from user access token
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

	// fmt.Printf("b: %v\n", b)

	var b RequestObject
	var handlerResp common_services.HandlerResponse
	if len(payload.RequestBody) != 0 {
		if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
			handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
			return handlerResp

		}
	}

	db := payload.Db

	valuesMap := make(map[string]interface{})

	query := `SELECT Q.* FROM(SELECT DISTINCT ON(s.space_id) space_id, s.name space_name, s.logo_url,s.type, CASE WHEN d.id IS NULL THEN false ELSE true END AS is_default,r.roles FROM spaces s inner join space_members mr ON mr.owner_space_id = s.space_id left join default_user_spaces d ON d.owner_space_id = s.space_id AND d.owner_user_id = mr.owner_user_id
	left join (select mr.owner_space_id,json_agg(json_build_object('r.name',r.name,'description',r.description,'id',r.id,'is_owner',r.is_owner)) as roles from member_roles mr inner join roles r on r.id=mr.role_id where mr.owner_user_id=@userID group by mr.owner_space_id)r on r.owner_space_id=s.space_id
	WHERE mr.owner_user_id = @userID`

	if len(b.SearchKeyword) > 0 {
		query += ` and s.name ilike @Keyword`
		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"
	}

	query += `) Q ORDER BY Q.is_default DESC, Q.space_name`
	valuesMap["userID"] = payload.UserID
	var spaceDetails []SpaceDetails

	// TODO
	res := db.Raw(query, valuesMap).Scan(&spaceDetails)

	if res.Error != nil {
		// RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	if res.RowsAffected < 1 {
		// RespondWithError(w, http.StatusNoContent, "NO RECORD FOUND!")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "NO RECORD FOUND!", Response{}, http.StatusNoContent)
		return handlerResp
	}

	var resp Response
	resp.Data = spaceDetails
	resp.Err = false
	resp.Msg = "Spaces listed successfully!"

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Spaces listed successfully!", resp, http.StatusOK)
	return handlerResp
}
