package check_entity_name

import (
	"encoding/json"
	"net/http"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {
	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db

	var nameExists Exists

	nameCheckQuery := `WITH 
    entity_space_mappings AS (
        SELECT * FROM entity_space_mappings esm WHERE esm.owner_space_id=?
    ),
    entities AS (
        SELECT * FROM entities as e 
        INNER JOIN entity_space_mappings esm ON e.entity_id = esm.owner_entity_id
    )
SELECT EXISTS(SELECT * FROM entities as et WHERE LOWER(et.label) = LOWER(?) and LOWER(et.type::varchar) = LOWER(?));`
	res := db.Raw(nameCheckQuery, b.SpaceID, b.Name, b.TypeID).Scan(&nameExists)

	if res.Error != nil {

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	if res.RowsAffected < 1 {

		handlerResp = common_services.BuildErrorResponse(true, "NO RECORD FOUND!", Response{}, http.StatusNoContent)
		return handlerResp
	}

	handlerResp = common_services.BuildResponse(true, "Data fetched successfully!", Response{Data: nameExists, Err: false, Msg: "Data fetched successfully!"}, http.StatusOK)
	return handlerResp
}
