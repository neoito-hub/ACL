package get_user_by_id

import (
	"encoding/json"
	"fmt"
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

	// fmt.Printf("b: %v\n", b)

	// TODO

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db

	valuesMap := make(map[string]interface{})

	usersQuery := `select u.email,u.user_id,u.user_name,u.full_name,u.phone,u.created_at as created_date,u.updated_at as updated_date from users u inner join space_members m on m.owner_user_id=u.user_id
	where u.user_id =@userID and m.owner_space_id =@spaceID`

	rolesQuery := `select json_agg(json_build_object('id',m.id,'name',r.name,'description',r.description,'is_owner',r.is_owner)) as roles from member_roles m left join roles r on r.id=m.role_id 
	where m.owner_user_id=@userID and m.owner_space_id =@spaceID `

	teamsQuery := `select json_agg(json_build_object('team_id',tm.id,'name',t.name,'description',t.description)) as teams from team_members tm left join teams t on t.team_id=tm.owner_team_id where tm.member_id=@userID and t.owner_id = @spaceID `

	query := fmt.Sprintf("select u.* from (select u.*,r.roles,t.teams from (%s)u left join (%s)r on true left join (%s)t on true)u", usersQuery, rolesQuery, teamsQuery)

	valuesMap["userID"] = b.UserID
	valuesMap["spaceID"] = b.SpaceID

	var user UserData

	db.Raw(query, valuesMap).Scan(&user)

	var resp Response
	resp.Data = user
	resp.Err = false
	resp.Msg = "User details retrieved successfully!"

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "User details retrieved successfully!", resp, http.StatusOK)
	return handlerResp
}
