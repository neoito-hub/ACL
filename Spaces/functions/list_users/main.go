package list_users

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

	// Validating and retreving user id from user access token
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

	// checking for permission
	// var blockauthor models.BlockAuthors

	// result := db.Model(&models.BlockAuthors{}).Where("user_id = ?", UserId).Find(&blockauthor)

	// fmt.Printf("block author is %v", blockauthor)

	// if result.RowsAffected == 0 {
	// 	var resp Response
	// 	resp.Data = "block author not found"
	// 	resp.Err = true
	// 	resp.Msg = "Block Author Not Found."

	// 	respondWithJSON(w, http.StatusNotFound, resp)
	// 	return
	// }

	SortColumns := make(map[string]string)

	SortColumns[b.Active] = "u." + b.Active
	SortColumns["createdAt"] = "u.created_date"
	SortColumns["updatedAt"] = "u.updated_date"

	var query string
	var usersFilterQuery string

	valuesMap := make(map[string]interface{})
	isDynamicQuery := false

	if b.IsKeywordSearch {
		isDynamicQuery = true
		usersFilterQuery = usersFilterQuery + " and (u.user_name  ilike @Keyword or u.email ilike @Keyword or u.full_name ilike @Keyword) "

		valuesMap["Keyword"] = "%" + b.Conditions.Keyword + "%"
	}

	usersQuery := fmt.Sprintf(`select u.email,MAX(CASE WHEN r.is_owner THEN 1 ELSE 0 END) as is_owner,u.user_id,u.user_name,u.full_name,u.phone,m.created_at as created_date,m.updated_at as updated_date from users u inner join space_members m on m.owner_user_id=u.user_id left join member_roles mr on mr.owner_user_id=u.user_id and mr.owner_space_id=m.owner_space_id left join roles r on r.owner_space_id=mr.owner_space_id and r.id=mr.role_id where m.owner_space_id=@spaceID %s
            group by u.email,u.user_id,u.user_name,u.full_name,u.phone,m.created_at,m.updated_at`, usersFilterQuery)

	rolesQuery := `select m.owner_user_id,json_agg(json_build_object('id',r.id,'name',r.name,'description',r.description,'is_owner',is_owner)) as roles from member_roles m left join roles r on r.id=m.role_id
	where m.owner_space_id=@spaceID group by m.owner_user_id`

	teamsQuery := `select tm.member_id,json_agg(json_build_object('team_id',t.team_id,'name',t.name,'description',t.description))  as teams from team_members tm inner join teams t on t.team_id=tm.owner_team_id where t.owner_id=@spaceID group by tm.member_id`

	countQuery := fmt.Sprintf(`select count(email) as total_count from (select  u.email from users u inner join space_members m on m.owner_user_id=u.user_id where m.owner_space_id=@spaceID %s) q`, usersFilterQuery)

	query = fmt.Sprintf("select u.* from (select u.*,r.roles,t.teams from (%s)u left join (%s)r on r.owner_user_id=u.user_id left join (%s)t on t.member_id=u.user_id)u", usersQuery, rolesQuery, teamsQuery)

	var orderByString string
	var limitString string
	var offsetString string

	valuesMap["spaceID"] = b.SpaceID

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
