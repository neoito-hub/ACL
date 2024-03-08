package list_roles

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

	var roles []RolesListData

	var filterString string

	valuesMap := make(map[string]interface{})
	rolesQuery1 := `WITH
    cte AS
    (
        SELECT
		mr.role_id, u.user_id, u.user_name,u.email, u.full_name,
		ROW_NUMBER() OVER (PARTITION BY mr.role_id ORDER BY u.email) AS rn
        FROM member_roles mr
		inner join users u on u.user_id = mr.owner_user_id
		where mr.owner_space_id =@owner_space_id
    ),
	
	members AS
	(
		SELECT role_id,
		JSONB_AGG(jsonb_build_object('user_id',user_id,'user_name',user_name,'email',email,'full_name', full_name) ORDER BY email) FILTER (WHERE rn <= @display_member_count) AS members
		FROM cte GROUP BY role_id
	)
	
	SELECT DISTINCT r.id as role_id, r.name, r.description, r.owner_space_id as space_id, (
		select COUNT(mr.role_id) from member_roles mr where mr.role_id = r.id
	) as member_count, m.members, r.is_owner FROM roles r 
	LEFT JOIN member_roles mr ON mr.role_id = r.id and mr.owner_user_id = @member_id
	LEFT JOIN (
		select mr.owner_space_id from member_roles mr
		inner join roles r on mr.role_id=r.id
		where mr.owner_space_id=@owner_space_id and mr.owner_user_id=@member_id and r.is_owner
		  ) o on o.owner_space_id = r.owner_space_id and mr.id is null
	LEFT JOIN members m ON m.role_id = r.id`
	filterString = ` where r.owner_space_id =@owner_space_id and (mr.id is not null or o.owner_space_id is not null)`
	rolesQuery2 := ` group by r.id`

	countQuery := `SELECT COUNT(id) as total_count FROM (SELECT r.id FROM roles r
		LEFT JOIN member_roles mr ON mr.role_id = r.id and mr.owner_user_id = @member_id
		LEFT JOIN (
			select mr.owner_space_id from member_roles mr
			inner join roles r on mr.role_id=r.id
			where mr.owner_space_id=@owner_space_id and mr.owner_user_id=@member_id and r.is_owner
			  ) o on o.owner_space_id = r.owner_space_id and mr.id is null`

	if len(b.SearchKeyword) > 0 {
		searchFilter := ""
		GenerateNonParameterisedQuery(&searchFilter, " r.name ilike @Keyword ", "and", filterString)

		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"

		AttachToMainfilter(&searchFilter, &filterString)

	}

	if b.DisplayMemberCount == 0 {
		b.DisplayMemberCount = 4
	}

	rolesQuery1 = rolesQuery1 + filterString + rolesQuery2 + ", m.members order by r.name limit @limit offset @offset"
	countQuery = countQuery + filterString + rolesQuery2 + ") q"
	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset

	valuesMap["member_id"] = userData.UserID
	valuesMap["owner_space_id"] = b.SpaceID
	valuesMap["display_member_count"] = b.DisplayMemberCount

	// get roles
	res := db.Raw(rolesQuery1, valuesMap).Scan(&roles)

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

	resultData.Data = roles

	// RespondWithJSON(w, http.StatusOK, Response{Data: resultData, Err: false, Msg: "Roles listed successfully!"})

	handlerResp = common_services.BuildResponse(false, "Roles listed successfully!", Response{Data: resultData, Err: false, Msg: "Roles listed successfully!"}, http.StatusOK)
	return handlerResp
}
