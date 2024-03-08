package user_list_available_entities

import (
	"encoding/json"
	"net/http"
	"strings"

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

	var entitiesList []EntitiesListData

	valuesMap := make(map[string]interface{})
	availableEntitiesQuery := `
	with acl_entities as (
		select et.entity_id,et.type,et.label,et.created_at from
		pol_grp_subs_entity_mappings etmap
		inner join entities et on et.entity_id=etmap.owner_entity_id
		inner join ac_pol_grp_subs pgs 
		on etmap.pol_grp_subs_id=pgs.id and pgs.owner_space_id=@owner_space_id
		left join ac_pol_grps polgrp on polgrp.id=pgs.ac_pol_grp_id
		
		LEFT JOIN (select mr.role_id, r.name,r.is_owner  from member_roles mr
			inner join roles r on r.id = mr.role_id
			where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id
		   ) r on r.role_id = pgs.role_id
 LEFT JOIN (select t.team_id, name from team_members tm 
			inner join teams t on tm.owner_team_id=t.team_id 
			where tm.member_id=@owner_user_id and t.owner_id=@owner_space_id
		   ) t ON t.team_id = pgs.owner_team_id
		   WHERE  polgrp.name in @policy_group_list and 
		  (pgs.owner_user_id=@owner_user_id OR r.role_id is not null or t.team_id is not null)
	),
	space_entities as (select distinct acl.type from acl_entities acl where acl.entity_id=acl.type::varchar),
owner_spaces as (select mr.role_id, r.name,r.is_owner,mr.owner_space_id  from member_roles mr
	   inner join roles r on r.id = mr.role_id
	   where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id and r.is_owner),
	owner_entities as (
		select et.entity_id,et.type,et.label,et.created_at from
		entities et
		left join entity_space_mappings esm on esm.owner_entity_id=et.entity_id and esm.owner_space_id=@owner_space_id
		left JOIN owner_spaces r on true
		left join space_entities sp on et.type=sp.type
		   WHERE  (r.is_owner or  sp.type is not null) and (esm.owner_entity_id is not null)
	)
    select et.* from 
	(SELECT acl_entities.* from acl_entities union 
	select owner_entities.* from owner_entities)et 
	 `

	availableEntitiesCountQuery := `
	with acl_entities as (
		select et.entity_id,et.type,et.label,et.created_at from
		pol_grp_subs_entity_mappings etmap
		inner join entities et on et.entity_id=etmap.owner_entity_id
		inner join ac_pol_grp_subs pgs 
		on etmap.pol_grp_subs_id=pgs.id and pgs.owner_space_id=@owner_space_id
		left join ac_pol_grps polgrp on polgrp.id=pgs.ac_pol_grp_id
		LEFT JOIN (select mr.role_id, r.name,r.is_owner  from member_roles mr
			inner join roles r on r.id = mr.role_id
			where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id
		   ) r on r.role_id = pgs.role_id
 LEFT JOIN (select t.team_id, name from team_members tm 
			inner join teams t on tm.owner_team_id=t.team_id 
			where tm.member_id=@owner_user_id and t.owner_id=@owner_space_id
		   ) t ON t.team_id = pgs.owner_team_id
		   WHERE  polgrp.name in @policy_group_list and 
		  (pgs.owner_user_id=@owner_user_id OR r.role_id is not null or t.team_id is not null)
	),
	space_entities as (select distinct acl.type from acl_entities acl where acl.entity_id=acl.type::varchar),
owner_spaces as (select mr.role_id, r.name,r.is_owner,mr.owner_space_id  from member_roles mr
	   inner join roles r on r.id = mr.role_id
	   where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id and r.is_owner),
	owner_entities as (
		select et.entity_id,et.type,et.label,et.created_at from
		entities et
		left join entity_space_mappings esm on esm.owner_entity_id=et.entity_id
		and esm.owner_space_id=@owner_space_id
		left JOIN owner_spaces r on r.owner_space_id=esm.owner_space_id
		left join space_entities sp on et.type=sp.type
		   WHERE  (r.is_owner or sp.type is not null) and (esm.owner_entity_id is not null)
	)

	select count(et.entity_id) from 
	(SELECT acl_entities.* from acl_entities union 
	select owner_entities.* from owner_entities)et 
	`

	filterString := "where et.type::varchar!=et.entity_id"

	if len(b.SearchKeyword) > 0 {
		searchFilter := ""

		//for the attached permissions
		GenerateNonParameterisedQuery(&searchFilter, " et.label ilike @Keyword  ", "and", filterString)

		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"

		AttachToMainfilter(&searchFilter, &filterString)

	}

	if len(b.EntityTypes) > 0 {
		entityTypeFilter := ""
		GenerateNonParameterisedQuery(&entityTypeFilter, " et.type in @entity_types ", "and", filterString)

		valuesMap["entity_types"] = b.EntityTypes

		AttachToMainfilter(&entityTypeFilter, &filterString)

	}

	SortColumns := make(map[string]string)
	SortDirections := make(map[string]string)
	if len(b.SortColumn) == 0 {
		b.SortColumn = "CreatedAt"
	}
	if len(b.SortDirection) == 0 {
		b.SortDirection = "desc"
	}

	SortColumns[b.SortColumn] = b.SortColumn
	SortDirections[b.SortDirection] = b.SortDirection

	SortColumns["CreatedAt"] = "et.created_at"
	SortColumns["UpdatedAt"] = "et.updated_at"
	SortDirections["desc"] = "desc"
	SortDirections["asc"] = "asc"

	orderByString := ` order by ` + SortColumns[b.SortColumn] + " " + SortDirections[strings.ToLower(b.SortDirection)]

	availableEntitiesQuery = availableEntitiesQuery + filterString + orderByString + " limit @limit offset @offset"
	availableEntitiesCountQuery = availableEntitiesCountQuery + filterString
	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset
	valuesMap["owner_user_id"] = payload.UserID
	valuesMap["owner_space_id"] = payload.SpaceID
	valuesMap["policy_group_list"] = []string{"ACL Access"}

	res := db.Raw(availableEntitiesQuery, valuesMap).Scan(&entitiesList)

	if res.Error != nil {

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	var resultData ResultData

	countRes := db.Raw(availableEntitiesCountQuery, valuesMap).Scan(&resultData.TotalCount)

	if countRes.Error != nil {

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	resultData.Data = entitiesList

	handlerResp = common_services.BuildResponse(false, "Available entities for user listed successfully!", Response{Data: resultData, Err: false, Msg: "Available entities for user listed successfully!"}, http.StatusOK)
	return handlerResp
}
