package update_user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

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
	// // shieldUser, shieldVerifyError := VerifyAndGetUser(w, r)
	// // if shieldVerifyError != nil {
	// // 	fmt.Printf("shieldVerifyError: %v\n", shieldVerifyError)
	// // 	RespondWithError(w, http.StatusUnauthorized, shieldVerifyError.Error())

	// // 	return
	// // }

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

	updateValue := make(map[string]interface{})

	fields := reflect.TypeOf(b.UserDetails)
	values := reflect.ValueOf(b.UserDetails)
	num := fields.NumField()

	for i := 0; i < num; i++ {
		fieldName := fields.Field(i).Name
		value := values.Field(i)

		var convertedValue interface{}

		convertIntBase := 10

		switch value.Kind() {
		case reflect.Bool:
			convertedValue = value.Bool()
		case reflect.Int:
			convertedValue = strconv.FormatInt(value.Int(), convertIntBase)
		case reflect.Interface:
			convertedValue = value.Interface()
		default:
			convertedValue = value.String()
		}

		if convertedValue != nil {
			updateValue[fieldName] = convertedValue
		}
	}

	var tx = db.Begin()

	// tx.Model(&models.User{}).Where("user_id = ?", b.UserID).Updates(updateValue)

	if len(b.DeletedRoleIds) > 0 {
		tx.Exec("delete from member_roles mr where mr.id in (?) and EXISTS (SELECT 1 FROM roles r WHERE r.id = mr.role_id AND r.is_owner IS NOT TRUE) ", b.DeletedRoleIds)
	}

	if len(b.DeletedTeamIds) > 0 {
		tx.Exec("delete from team_members tm where tm.id in (?) ", b.DeletedTeamIds)
	}

	valuesMap := make(map[string]interface{})

	usersQuery := `select u.email,u.user_id,u.user_name,u.full_name,u.phone,u.created_at as created_date,u.updated_at as updated_date from users u inner join member_roles m on m.owner_user_id=u.user_id
	where u.user_id =@userID`

	rolesQuery := `select json_agg(json_build_object('id',m.id,'name',r.name,'description',r.description,'is_owner',r.is_owner)) as roles from member_roles m left join roles r on r.id=m.role_id 
	where m.owner_user_id=@userID `

	teamsQuery := `select json_agg(json_build_object('team_id',tm.id,'name',t.name,'description',t.description)) as teams from team_members tm left join teams t on t.team_id=tm.owner_team_id where tm.member_id=@userID `

	query := fmt.Sprintf("select u.* from (select u.*,r.roles,t.teams from (%s)u left join (%s)r on true left join (%s)t on true)u", usersQuery, rolesQuery, teamsQuery)

	valuesMap["userID"] = b.UserID

	var user UserData

	db.Raw(query, valuesMap).Scan(&user)

	tx.Commit()
	resp := Response{
		Data: user,
		Err:  false,
		Msg:  "Member details updated successfully!",
	}

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Member details updated successfully!", resp, http.StatusOK)
	return handlerResp
}
