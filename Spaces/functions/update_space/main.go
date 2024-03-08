package update_space

import (
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/neoito-hub/ACL-Block/Data-Models/models"
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

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")

	if !re.MatchString(b.Name) || !re.MatchString(b.BusinessName) {
		handlerResp = common_services.BuildErrorResponse(true, "No special characters allowed other than underscore for name and business name", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	db := payload.Db

	shieldUser := ShieldUserData{
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}

	// Get existing space data
	var spaceData SpaceDetails

	db.Raw("SELECT s.space_id, s.name, s.type, s.email, s.country, s.business_name, s.address, s.business_category, s.description, s.market_place_id, s.developer_portal_access, s.meta_data, s.logo_url FROM spaces s INNER JOIN member_roles mr ON mr.owner_space_id = s.space_id INNER JOIN roles r ON r.owner_space_id = s.space_id AND r.id = mr.role_id WHERE s.space_id = ? AND mr.owner_user_id = ? AND r.is_owner = true", b.SpaceID, shieldUser.UserID).Scan(&spaceData)

	updateValue := make(map[string]interface{})
	// data, mErr := json.Marshal(spaceData)

	// if mErr != nil {
	// 	RespondWithError(w, http.StatusBadRequest, "Something went wrong")
	// 	return
	// }

	// jsonErr := json.Unmarshal(data, &updateValue)

	// if jsonErr != nil {
	// 	RespondWithError(w, http.StatusBadRequest, "Something went wrong")
	// 	return
	// }

	if b.Type != "" && b.Type != spaceData.Type {
		// RespondWithError(w, http.StatusBadRequest, "Change type of space is not allowed")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Change type of space is not allowed", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	if spaceData.Type == "P" && b.Country != "" && b.BusinessCategory != "" && b.BusinessName != "" && b.Address != "" {
		// RespondWithError(w, http.StatusBadRequest, "Change business data of space is not allowed for personal space")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Change business data of space is not allowed for personal space", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	// Update existing data object with new data
	fields := reflect.TypeOf(b)
	values := reflect.ValueOf(b)
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
		case reflect.String:
			if len(value.String()) > 0 {
				convertedValue = value.String()
			} else {
				convertedValue = nil
			}

		case reflect.Interface:
			convertedValue = value.Interface()
		default:
			convertedValue = value.String()
		}

		if convertedValue != nil {
			updateValue[fieldName] = convertedValue
		}
	}

	// Update space data
	var spaceDetails SpaceDetails

	res := db.Model(&models.Space{}).Where("space_id = ?", b.SpaceID).Updates(updateValue).Scan(&spaceDetails)

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

	resp := Response{
		Data: spaceDetails,
		Err:  false,
		Msg:  "Space details updated successfully!",
	}

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Space details updated successfully!", resp, http.StatusOK)
	return handlerResp
}
