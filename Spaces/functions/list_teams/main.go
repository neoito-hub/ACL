package list_teams

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

	var teams []TeamsListData

	var filterString string

	valuesMap := make(map[string]interface{})
	teamsQuery1 := `WITH
    cte AS
    (
        SELECT
		tm.owner_team_id, u.user_id, u.user_name,u.email, u.full_name,
		ROW_NUMBER() OVER (PARTITION BY tm.owner_team_id ORDER BY u.email) AS rn
        FROM teams t
      	inner join team_members tm on tm.owner_team_id = t.team_id
		inner join users u on u.user_id = tm.member_id
		where t.owner_id = @owner_id
    ),
	
	members AS
	(
		SELECT owner_team_id,
		JSONB_AGG(jsonb_build_object('user_id',user_id,'user_name',user_name,'email',email,'full_name', full_name) ORDER BY email) FILTER (WHERE rn <= @display_member_count) AS members
		FROM cte GROUP BY owner_team_id
	)

	SELECT DISTINCT t.team_id as team_id, t.name, t.description, t.owner_id as space_id, (
		select COUNT(tm.id) from team_members tm where tm.owner_team_id = t.team_id
	) as member_count, m.members, tm.is_owner FROM teams t
	LEFT JOIN team_members tm ON tm.owner_team_id = t.team_id AND tm.member_id=@member_id
	LEFT JOIN (select mr.owner_space_id from member_roles mr
		inner join roles r on mr.role_id=r.id
		where mr.owner_space_id=@owner_id and mr.owner_user_id=@member_id and r.is_owner
	   ) o on o.owner_space_id = t.owner_id and tm.id is null
	LEFT JOIN members m ON m.owner_team_id = t.team_id`

	filterString = ` where t.owner_id =@owner_id AND (tm.id is not null OR o.owner_space_id is not null)`
	teamsQuery2 := ` group by t.team_id, tm.is_owner`

	countQuery := `SELECT COUNT(team_id) as total_count FROM (SELECT t.team_id FROM teams t 
		LEFT JOIN team_members tm ON tm.owner_team_id = t.team_id AND tm.member_id=@member_id
		LEFT JOIN (select mr.owner_space_id from member_roles mr
			inner join roles r on mr.role_id=r.id
			where mr.owner_space_id=@owner_id and mr.owner_user_id=@member_id and r.is_owner
		   ) o on o.owner_space_id = t.owner_id and tm.id is null`

	if len(b.SearchKeyword) > 0 {
		searchFilter := ""
		GenerateNonParameterisedQuery(&searchFilter, " t.name ilike @Keyword ", "and", filterString)

		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"

		AttachToMainfilter(&searchFilter, &filterString)

	}

	if b.DisplayMemberCount == 0 {
		b.DisplayMemberCount = 4
	}

	teamsQuery1 = teamsQuery1 + filterString + teamsQuery2 + ", m.members order by t.name limit @limit offset @offset"
	countQuery = countQuery + filterString + teamsQuery2 + ") q"
	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset

	valuesMap["member_id"] = userData.UserID
	valuesMap["owner_id"] = b.SpaceID
	valuesMap["display_member_count"] = b.DisplayMemberCount

	// get teams
	res := db.Raw(teamsQuery1, valuesMap).Scan(&teams)

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

	resultData.Data = teams

	// RespondWithJSON(w, http.StatusOK, Response{Data: resultData, Err: false, Msg: "Teams listed successfully!"})

	handlerResp = common_services.BuildResponse(false, "Teams listed successfully!", Response{Data: resultData, Err: false, Msg: "Teams listed successfully!"}, http.StatusOK)
	return handlerResp
}
