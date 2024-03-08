package user_list_existing_pol_grp_subs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	userQuery1 := `SELECT pg.name, pg.description,pg.is_predefined,pgs.permission_id,json_agg(json_build_object('role_id',r.role_id,'role_name',r.name)) filter (where r.role_id is not null) roles, json_agg(json_build_object('team_id',t.team_id,'team_name',t.name))  filter (where t.team_id is not null) teams,  min(pgs.created_at) created_at, max(pgs.updated_at) updated_at FROM (select * from ac_pol_grp_subs) pgs 
	INNER JOIN ac_pol_grps pg ON pg.id = pgs.ac_pol_grp_id
	LEFT JOIN (select mr.role_id, r.name  from member_roles mr
			   inner join roles r on r.id = mr.role_id
			   where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id
			  ) r on r.role_id = pgs.role_id
	LEFT JOIN (select t.team_id, name from team_members tm 
			   inner join teams t on tm.owner_team_id=t.team_id 
			   where tm.member_id=@owner_user_id and t.owner_id=@owner_space_id
			  ) t ON t.team_id = pgs.owner_team_id
	
			  `

	filterString = ` WHERE  pgs.owner_space_id=@owner_space_id AND 
	(owner_user_id=@owner_user_id OR r.role_id is not null or t.team_id is not null)`

	countQuery := `SELECT COUNT(DISTINCT concat(pg.id,pgs.permission_id)) as total_count FROM (select * from ac_pol_grp_subs) pgs
	INNER JOIN ac_pol_grps pg ON pg.id = pgs.ac_pol_grp_id
	LEFT JOIN (
		select mr.role_id  from member_roles mr where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id
	) r on r.role_id = pgs.role_id
	LEFT JOIN (
		select t.team_id, name from team_members tm
		inner join teams t on tm.owner_team_id=t.team_id
		where tm.member_id=@owner_user_id and t.owner_id=@owner_space_id
	) t ON t.team_id = pgs.owner_team_id`

	if len(b.SearchKeyword) > 0 {
		searchFilter := ""
		GenerateNonParameterisedQuery(&searchFilter, " pg.name ilike @Keyword ", "and", filterString)

		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"

		AttachToMainfilter(&searchFilter, &filterString)

	}

	SortColumns := make(map[string]string)
	SortDirections := make(map[string]string)
	if len(b.SortColumn) == 0 {
		b.SortColumn = "updatedAt"
	}
	if len(b.SortDirection) == 0 {
		b.SortDirection = "desc"
	}

	SortColumns[b.SortColumn] = b.SortColumn
	SortDirections[b.SortDirection] = b.SortDirection

	SortColumns["createdAt"] = "created_at"
	SortColumns["updatedAt"] = "updated_at"
	SortDirections["desc"] = "desc"
	SortDirections["asc"] = "asc"

	orderByString := ` order by pg.` + SortColumns[b.SortColumn] + " " + SortDirections[strings.ToLower(b.SortDirection)]

	userQuery1 = fmt.Sprintf(`select pg.* from (%s)pg`, userQuery1+filterString+"group by pg.name, pg.description,pg.is_predefined,pgs.permission_id") +
		orderByString + " limit @limit offset @offset"
	countQuery = countQuery + filterString
	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset

	valuesMap["owner_user_id"] = b.UserID
	valuesMap["owner_space_id"] = b.SpaceID

	res := db.Raw(userQuery1, valuesMap).Scan(&policiesList)

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

	// RespondWithJSON(w, http.StatusOK, Response{Data: resultData, Err: false, Msg: "Existing policy subscription for user listed successfully!"})

	handlerResp = common_services.BuildResponse(false, "Existing policy subscription for user listed successfully!", Response{Data: resultData, Err: false, Msg: "Existing policy subscription for user listed successfully!"}, http.StatusOK)
	return handlerResp
}
