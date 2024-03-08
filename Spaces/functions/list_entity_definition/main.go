package list_entity_definition

import (
	"encoding/json"
	"net/http"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
	"gorm.io/gorm"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {

	var b RequestObject
	var handlerResp common_services.HandlerResponse
	if len(payload.RequestBody) != 0 {
		if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
			handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
			return handlerResp
		}
	}

	db := payload.Db

	var res *gorm.DB

	valuesMap := make(map[string]interface{})

	query := `select ed.id , ed.display_name, ed.name 
	from entity_type_definitions as ed `

	if len(b.SearchKeyword) > 0 {
		query += `where ed.display_name ilike @Keyword`
		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"
	}

	query += ` order by ed.display_name DESC `
	var entityDefinitionDetails []entityDefinitionDetails

	if len(b.SearchKeyword) > 0 {
		res = db.Raw(query, valuesMap).Scan(&entityDefinitionDetails)
	} else {
		res = db.Raw(query).Scan(&entityDefinitionDetails)
	}

	if res.Error != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	if res.RowsAffected < 1 {
		handlerResp = common_services.BuildErrorResponse(true, "NO RECORD FOUND!", Response{}, http.StatusNoContent)
		return handlerResp
	}

	var resp Response
	resp.Data = entityDefinitionDetails
	resp.Err = false
	resp.Msg = "Entity Definitions listed successfully!"

	handlerResp = common_services.BuildResponse(false, "Entity Definitions listed successfully!", resp, http.StatusOK)
	return handlerResp
}
