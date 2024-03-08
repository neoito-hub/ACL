package list_invited_users

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

	SortColumns[b.Active] = "i." + b.Active
	SortColumns["createdAt"] = "i.created_at"
	SortColumns["updatedAt"] = "i.updated_at"

	var query string
	var inviteFilterQuery string

	valuesMap := make(map[string]interface{})
	isDynamicQuery := false

	if b.IsKeywordSearch {
		isDynamicQuery = true

		inviteFilterQuery = inviteFilterQuery + " and i.email ilike @Keyword "

		valuesMap["Keyword"] = "%" + b.Conditions.Keyword + "%"
	}

	query = fmt.Sprintf(`select i.email,i.id as invite_id,i.created_at as created_date,i.updated_at as updated_date, CASE WHEN i.expires_at < now() THEN TRUE ELSE FALSE END as expired, i.expires_at, i.status,json_agg(json_build_object('role_id',r.id,'name',r.name,'description',r.description,'invite_details_id',id.id)) filter (where r.id is not null) as roles,json_agg(json_build_object('team_id',t.team_id,'name', t.name,'description',t.description,'invite_details_id',id.id)) filter (where t.team_id is not null) as teams from invites i left join invite_details id on id.invite_id=i.id left join teams t on t.team_id=id.invited_team_id left join roles r on r.id = id.invited_role_id left join (select u.email from users u inner join space_members m on m.owner_user_id=u.user_id where m.owner_space_id=@spaceID) u on u.email = i.email where u.email is null and id.invited_space_id=@spaceID and i.invite_type=1 and i.status != 2 %s  group by i.email,i.id,i.created_at,i.updated_at, i.expires_at, i.status`, inviteFilterQuery)

	countQuery := fmt.Sprintf(`select count(email) as total_count from (select i.email from invites i left join invite_details id on id.invite_id=i.id left join (select u.email from users u inner join space_members m on m.owner_user_id=u.user_id where m.owner_space_id=@spaceID) u on u.email = i.email where u.email is null and id.invited_space_id=@spaceID and i.invite_type=1 and i.status != 2 %s group by i.email,i.id,i.created_at,i.updated_at, i.expires_at, i.status) q`, inviteFilterQuery)

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
	resp.Msg = "Invite list fetched successfully!"

	// respondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Invite list fetched successfully!", resp, http.StatusOK)
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
