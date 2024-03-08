package list_spaces_detailed

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

	stateCondition := ``
	if b.State == 1 {
		stateCondition = ` and s.type='P'`
	} else if b.State == 2 {
		stateCondition = ` and s.type='B'`
	}

	query := `SELECT Q.* FROM(
		SELECT s.space_id, s.name space_name, s.logo_url,s.type, s.created_at,
		(select count(owner_user_id) from space_members where owner_space_id = s.space_id) as member_count,
		(select count(owner_entity_id) from entity_space_mappings where owner_space_id = s.space_id) as entity_count,
		CASE WHEN mr.owner_user_id = @userID THEN TRUE ELSE FALSE END as my_space,
		mr.owner_user_id, mr.full_name, mr.email, mr.user_name
		 FROM spaces s 
		inner join space_members sm ON sm.owner_space_id = s.space_id 
		left join (
			select distinct on (mr.owner_space_id) mr.owner_space_id, mr.owner_user_id, u.full_name, u.email, u.user_name from member_roles mr 
			inner join roles r on mr.role_id=r.id 
			inner join users u on u.user_id=mr.owner_user_id
			where r.is_owner order by mr.owner_space_id, mr.created_at
		) mr on mr.owner_space_id = s.space_id
		WHERE sm.owner_user_id = @userID` + stateCondition

	countQuery := `SELECT count(s.space_id) as total_count
	FROM spaces s 
   inner join space_members sm ON sm.owner_space_id = s.space_id 
   left join (
	   select distinct on (mr.owner_space_id) mr.owner_space_id, mr.owner_user_id, u.full_name, u.email, u.user_name from member_roles mr 
	   inner join roles r on mr.role_id=r.id 
	   inner join users u on u.user_id=mr.owner_user_id
	   where r.is_owner order by mr.owner_space_id, mr.created_at
   ) mr on mr.owner_space_id = s.space_id
   WHERE sm.owner_user_id = @userID` + stateCondition

	if len(b.SearchKeyword) > 0 {
		query += ` and s.name ilike @Keyword`
		countQuery += ` and s.name ilike @Keyword`
		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"
	}

	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset

	query += ` ORDER BY s.name limit @limit offset @offset) Q`
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

	var resultData ResultData

	countRes := db.Raw(countQuery, valuesMap).Scan(&resultData.TotalCount)

	if countRes.Error != nil {
		// RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	resultData.Data = spaceDetails

	var resp Response
	resp.Data = resultData
	resp.Err = false
	resp.Msg = "Spaces listed successfully!"

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Spaces listed successfully!", resp, http.StatusOK)
	return handlerResp
}
