package teams_list_permissions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db

	var policiesList []PoliciesListData

	var filterString string

	valuesMap := make(map[string]interface{})
	existingPermissionsQuery := `
	with pg_subs as (
		select pgs.id,pg.display_name,pgs.permission_id,pgs.owner_team_id,pgs.owner_space_id,unnest(pg.entity_types) as entity_type from
		ac_pol_grp_subs pgs 
		INNER JOIN ac_pol_grps pg ON pg.id = pgs.ac_pol_grp_id
 		INNER JOIN teams t ON t.team_id = pgs.owner_team_id 
		WHERE  pgs.owner_space_id=@owner_space_id AND t.team_id=@owner_team_id
	)
	
SELECT acper.id as permission_id,acper.display_name as name, acper.description,json_agg(distinct jsonb_build_object('subs_id',pg_subs.id,'policy_group_name',pg_subs.display_name)) policy_groups, count(pg_subs.id) as pg_count,jsonb_agg(distinct jsonb_build_object('entity_id',et.entity_id,'entity_type',et.type,'label',et.label)) filter (where et.entity_id is not null) as entities,json_object_agg(distinct coalesce(pg_subs.entity_type,0),true) as entity_types FROM 
	pg_subs
	inner join ac_permissions acper on acper.id=pg_subs.permission_id
	left join pol_grp_subs_entity_mappings etmap on etmap.pol_grp_subs_id=pg_subs.id
	left join entities et on et.entity_id=etmap.owner_entity_id

	 `

	filterString = ``

	existingPermissionsCountQuery := `
	with pg_subs as (
		select pgs.id,pg.display_name,pgs.permission_id,pgs.owner_team_id,pgs.owner_space_id from
		ac_pol_grp_subs pgs 
		INNER JOIN ac_pol_grps pg ON pg.id = pgs.ac_pol_grp_id
 		INNER JOIN teams t ON t.team_id = pgs.owner_team_id 
		WHERE  pgs.owner_space_id=@owner_space_id AND t.team_id=@owner_team_id
	)
	SELECT COUNT(*) as total_count FROM 
	(select distinct acper.id from 
	pg_subs
	inner join ac_permissions acper on acper.id=pg_subs.permission_id filterstring)acper
	`

	if len(b.SearchKeyword) > 0 {
		//for the attached permissions
		searchFilter := ""
		GenerateNonParameterisedQuery(&searchFilter, " acper.display_name ilike @Keyword ", "and", filterString)

		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"

		AttachToMainfilter(&searchFilter, &filterString)

	}

	SortColumns := make(map[string]string)
	SortDirections := make(map[string]string)
	if len(b.SortColumn) == 0 {
		b.SortColumn = "PolicyCount"
	}
	if len(b.SortDirection) == 0 {
		b.SortDirection = "desc"
	}

	SortColumns[b.SortColumn] = b.SortColumn
	SortDirections[b.SortDirection] = b.SortDirection

	SortColumns["PolicyCount"] = "pg_count"
	SortColumns["updatedAt"] = "updated_at"
	SortDirections["desc"] = "desc"
	SortDirections["asc"] = "asc"

	orderByString := ` order by ` + SortColumns[b.SortColumn] + " " + SortDirections[strings.ToLower(b.SortDirection)]

	existingPermissionsQuery = existingPermissionsQuery + filterString + " group by acper.id " + orderByString + " limit @limit offset @offset"
	existingPermissionsCountQuery = strings.Replace(existingPermissionsCountQuery, "filterstring", filterString, 1)
	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset
	valuesMap["owner_team_id"] = b.TeamID
	valuesMap["owner_space_id"] = payload.SpaceID

	res := db.Raw(existingPermissionsQuery, valuesMap).Scan(&policiesList)

	if res.Error != nil {
		// RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	// if res.RowsAffected < 1 {
	// 	// RespondWithJSON(w, http.StatusNoContent, Response{Err: false, Msg: "NO RECORD FOUND!"})
	// 	// return

	// 	handlerResp = common_services.BuildErrorResponse(true, "NO RECORD FOUND!", Response{}, http.StatusNoContent)
	// 	return handlerResp
	// }

	var resultData ResultData

	countRes := db.Raw(existingPermissionsCountQuery, valuesMap).Scan(&resultData.TotalCount)

	if countRes.Error != nil {
		// RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	for k, permission := range policiesList {
		var convertedEntities []Entity
		var attachedEntities AttachedEntities
		spaceAccessEntities := make(Entities)
		addedEntities := make(Entities)

		json.Unmarshal(permission.Entities, &convertedEntities)

		for _, entity := range convertedEntities {
			if entity.EntityID == fmt.Sprintf("%v", entity.EntityType) {

				spaceAccessEntities[entity.EntityType] = append(spaceAccessEntities[entity.EntityType], entity)

			} else {

				addedEntities[entity.EntityType] = append(addedEntities[entity.EntityType], entity)
			}

		}

		attachedEntities.SpaceAccessEntities = spaceAccessEntities
		attachedEntities.AddedEntities = addedEntities

		policiesList[k].AttachedEntities = attachedEntities

	}

	resultData.Data = policiesList

	// RespondWithJSON(w, http.StatusOK, Response{Data: resultData, Err: false, Msg: "Existing policy subscription for user listed successfully!"})

	handlerResp = common_services.BuildResponse(false, "Existing permission subscription for team listed successfully!", Response{Data: resultData, Err: false, Msg: "Existing permission subscription for team listed successfully!"}, http.StatusOK)
	return handlerResp
}
