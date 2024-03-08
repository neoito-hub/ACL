package roles_list_users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

// Handler func
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

	// //closing connection to db
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

	SortColumns := make(map[string]string)

	SortColumns[b.Active] = "u." + b.Active
	SortColumns["createdAt"] = "u.created_date"
	SortColumns["updatedAt"] = "u.updated_date"

	var query string
	var usersFilterQuery string
	// var inviteFilterQuery string

	valuesMap := make(map[string]interface{})
	isDynamicQuery := false

	if b.IsKeywordSearch {
		isDynamicQuery = true
		usersFilterQuery = usersFilterQuery + " and (u.user_name  ilike @Keyword or u.email ilike @Keyword or u.full_name ilike @Keyword) "

		// inviteFilterQuery = inviteFilterQuery + " and (id.email ilike @Keyword) "

		valuesMap["Keyword"] = "%" + b.Conditions.Keyword + "%"
	}

	usersQuery := fmt.Sprintf(`select m.id as member_role_id, m.owner_space_id as space_id, u.email,0 as type,u.user_id,u.user_name,u.full_name,u.phone,m.created_at as created_date,m.updated_at as updated_date from users u inner join member_roles m on m.owner_user_id=u.user_id where m.role_id=@roleID %s `, usersFilterQuery)

	teamsQuery := `select m.member_id, t.owner_id as space_id, json_agg(json_build_object('id',t.team_id,'name',t.name,'description',t.description,'is_owner', m.is_owner)) as teams from team_members m inner join teams t on t.team_id=m.owner_team_id group by m.member_id, t.owner_id`

	// invitesQuery := fmt.Sprintf(`select null as member_role_id, email, 1 as type, null as user_id,null as user_name,null as full_name,null as phone,null as created_date,null as updated_date, null as teams from (select id.email from invite_details id inner join invites i on id.invite_id=i.id inner join roles r on r.id=id.invited_role_id where now()<i.expires_at and i.invite_type=1 and id.invited_role_id = @roleID and i.status=1 %s group by id.email) temp group by temp.email`, inviteFilterQuery)

	countQuery := fmt.Sprintf(`select count(email) as total_count from (select  u.email from users u inner join member_roles m on m.owner_user_id=u.user_id where m.role_id=@roleID %s) q`, usersFilterQuery)

	query = fmt.Sprintf("select u.* from (select u.member_role_id,u.email,u.type,u.user_id,u.user_name,u.full_name,u.phone,u.created_date,u.updated_date,t.teams from (%s)u left join (%s)t on t.member_id=u.user_id and t.space_id=u.space_id)u", usersQuery, teamsQuery)

	var orderByString string
	var limitString string
	var offsetString string

	valuesMap["roleID"] = b.RoleID

	orderByString = ` order by ` + SortColumns[b.Active] + " " + b.Direction

	//adding limit
	if b.PageLimit > 0 {
		limitString = ` LIMIT @pageLimit`
		valuesMap["pageLimit"] = b.PageLimit
	}

	//adding offest for which order by is mandatory
	if b.Offset >= 0 {
		// orderByString = ` order by @active`
		// valuesMap["active"] = SortColumns[b.Active] + " " + b.Direction

		offsetString = ` OFFSET @offset`
		valuesMap["offset"] = b.Offset
	}

	query += orderByString + limitString + offsetString

	var users []Data
	var resultData ResultData

	db.Raw(countQuery, valuesMap).Scan(&resultData.TotalCount)

	if isDynamicQuery {
		db.Raw(query, valuesMap).Scan(&users)
	} else {
		db.Raw(query, valuesMap).Scan(&users)
	}

	resultData.Data = users

	var resp Response
	resp.Data = resultData
	resp.Err = false
	resp.Msg = "Users list fetched successfully!"

	// respondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Users list fetched successfully!", resp, http.StatusOK)
	return handlerResp

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	var resp ErrResponse
	resp.Err = true
	resp.Msg = message
	respondWithJSON(w, code, resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dataPayload, err := json.Marshal(payload)
	if err != nil {
		code = http.StatusInternalServerError
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dataPayload)
}
