package roles_list_to_add_pol_grp_subs

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

	var policiesList []PoliciesListData

	var filterString string

	valuesMap := make(map[string]interface{})
	rolesQuery1 := `SELECT pg.id AS ac_pol_grp_id, pg.name, pg.description, pg.is_predefined, pg.created_at, pg.updated_at, pgs.id subs_id FROM ac_pol_grps pg LEFT JOIN (select * from ac_pol_grp_subs where permission_id is null) pgs ON pgs.ac_pol_grp_id=pg.id AND pgs.owner_space_id = @owner_space_id AND pgs.role_id = @role_id `
	filterString = ` WHERE  (pg.owner_space_id=@owner_space_id OR coalesce(pg.owner_space_id, '') =  '')`

	countQuery := `SELECT COUNT(pg.id) as total_count FROM ac_pol_grps pg LEFT JOIN ac_pol_grp_subs pgs ON pgs.ac_pol_grp_id=pg.id AND pgs.owner_space_id = @owner_space_id AND pgs.role_id = @role_id`

	if len(b.SearchKeyword) > 0 {
		searchFilter := ""
		GenerateNonParameterisedQuery(&searchFilter, " pg.name ilike @Keyword ", "and", filterString)

		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"

		AttachToMainfilter(&searchFilter, &filterString)

	}

	rolesQuery1 = rolesQuery1 + filterString + " limit @limit offset @offset"
	countQuery = countQuery + filterString
	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset

	valuesMap["role_id"] = b.RoleID
	valuesMap["owner_space_id"] = b.SpaceID

	res := db.Raw(rolesQuery1, valuesMap).Scan(&policiesList)

	if res.Error != nil {
		// RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	if res.RowsAffected < 1 {
		// RespondWithJSON(w, http.StatusNoContent, Response{Err: false, Msg: "NO RECORD FOUND!"})
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

	resultData.Data = policiesList

	// RespondWithJSON(w, http.StatusOK, Response{Data: resultData, Err: false, Msg: "Policy subscription for roles listed successfully!"})

	handlerResp = common_services.BuildResponse(false, "Policy subscription for roles listed successfully!", Response{Data: resultData, Err: false, Msg: "Policy subscription for roles listed successfully!"}, http.StatusOK)
	return handlerResp
}
