package teams_list_entities

import (
	"encoding/json"
	"fmt"
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

	db := payload.Db

	//closing connection to db
	// sqlDB, dberr := db.DB()
	// if dberr != nil {
	// 	log.Fatalln(dberr)
	// }
	// defer sqlDB.Close()

	valuesMap := make(map[string]interface{})
	var entitiesData EntityData
	var limitString string
	var offsetString string
	var orderByString string
	var blocksFilterString string
	policyGroupName := `ACL Access`

	blocksFilterString = `where (set.entity_id is not null or oe.entity_id is not null) and et.type=@typeID and et.type::varchar!=et.entity_id `
	valuesMap["userID"] = payload.UserID
	valuesMap["spaceID"] = payload.SpaceID
	valuesMap["teamID"] = b.TeamID
	valuesMap["typeID"] = b.TypeID
	SortColumns := make(map[string]string)

	SortColumns[b.SortColumn] = "et." + b.SortColumn
	SortColumns["createdAt"] = "et.created_at"
	SortColumns["updatedAt"] = "et.updated_at"
	SortDirections := map[string]string{"desc": "desc", "asc": "asc"}

	if len(b.SortColumn) == 0 {
		b.SortColumn = "updatedAt"
	}
	if len(b.SortDirection) == 0 {
		b.SortColumn = "desc"
	}

	orderByString = ` order by ` + SortColumns[b.SortColumn] + " " + SortDirections[b.SortDirection]
	//adding limit
	if b.Limit > 0 {
		limitString = ` LIMIT @limit`
		valuesMap["limit"] = b.Limit
	}

	//adding offest for which order by is mandatory
	if b.Offset >= 0 {

		offsetString = ` OFFSET @offset`
		valuesMap["offset"] = b.Offset
	}

	if b.Conditions.IsKeywordSearch {
		GenerateNonParameterisedQuery(&blocksFilterString, "(et.label ilike @searchKeyword)", "and")
		valuesMap["searchKeyword"] = "%" + b.Conditions.Keyword + "%"
	}

	countQuery := fmt.Sprintf(`with subs_entities as (
		select etmap.owner_entity_id as entity_id from (select * from ac_pol_grp_subs) subs 
		inner join ac_pol_grps polgrp on polgrp.id=subs.ac_pol_grp_id and polgrp.entity_types @> array_append(array[]::integer[],@typeID) and polgrp.name='%s'
		inner join pol_grp_subs_entity_mappings etmap on etmap.pol_grp_subs_id=subs.id
		--joining with user roles
		left join (select mr.role_id from member_roles mr where mr.owner_user_id=@userID
			and mr.owner_space_id=@spaceID) ur on ur.role_id=subs.role_id 
		
		--joining with user teams
		left join (select t.team_id from team_members tm inner join teams t on tm.owner_team_id=t.team_id where tm.member_id=@userID
			and t.owner_id=@spaceID) ut on ut.team_id=subs.owner_team_id 
		
		where (subs.owner_space_id=@spaceID and subs.owner_user_id=@userID) or (ur.role_id is not null or ut.team_id is not null)
			),
		
		
			space_entities as (select distinct et.type from entities et inner join subs_entities sub on sub.entity_id=et.entity_id where et.entity_id=et.type::varchar),
			owner_spaces as (select mr.role_id, r.name,r.is_owner,mr.owner_space_id  from member_roles mr
				inner join roles r on r.id = mr.role_id
				where mr.owner_user_id=@userID and mr.owner_space_id=@spaceID and r.is_owner),
		
		--getting owner_entities for the admin acquired via space owner privileges
		owner_entities as (	select et.entity_id from
			entities et
			left join entity_space_mappings esm on esm.owner_entity_id=et.entity_id and esm.owner_space_id=@spaceID 
			LEFT JOIN owner_spaces r on r.owner_space_id=esm.owner_space_id
			left join space_entities sp on sp.type=et.type
			   WHERE  (r.is_owner or  sp.type is not null) and (esm.owner_entity_id is not null))
		
		
		select count(*) from 
					(select et.entity_id from (select et.* from entities et where et.type=@typeID) et  left join owner_entities oe on oe.entity_id=et.entity_id left join subs_entities set on set.entity_id=et.entity_id %s
					group by et.entity_id)et 
		
		
		`, policyGroupName, blocksFilterString)

	query := fmt.Sprintf(`with

--getting entities attached via acl for the admin
	subs_entities as (
select etmap.owner_entity_id as entity_id,etmap.id as entity_mapping_id from (select * from ac_pol_grp_subs) subs 
inner join ac_pol_grps polgrp on polgrp.id=subs.ac_pol_grp_id and polgrp.entity_types @> array_append(array[]::integer[],@typeID) and polgrp.name='%s'
inner join pol_grp_subs_entity_mappings etmap on etmap.pol_grp_subs_id=subs.id
--joining with user roles
left join (select mr.role_id from member_roles mr where mr.owner_user_id=@userID
	and mr.owner_space_id=@spaceID) ur on ur.role_id=subs.role_id 

--joining with user teams
left join (select t.team_id from team_members tm inner join teams t on tm.owner_team_id=t.team_id where tm.member_id=@userID
	and t.owner_id=@spaceID) ut on ut.team_id=subs.owner_team_id 

where (subs.owner_space_id=@spaceID and subs.owner_user_id=@userID) or (ur.role_id is not null or ut.team_id is not null)
	),
	space_entities as (select distinct et.type from entities et inner join subs_entities sub on sub.entity_id=et.entity_id where et.entity_id=et.type::varchar),
			owner_spaces as (select mr.role_id, r.name,r.is_owner,mr.owner_space_id  from member_roles mr
				inner join roles r on r.id = mr.role_id
				where mr.owner_user_id=@userID and mr.owner_space_id=@spaceID and r.is_owner),
		
		--getting owner_entities for the admin acquired via space owner privileges
		owner_entities as (	select et.entity_id from
			entities et
			left join entity_space_mappings esm on esm.owner_entity_id=et.entity_id and esm.owner_space_id=@spaceID 
			LEFT JOIN owner_spaces r on r.owner_space_id=esm.owner_space_id
			left join space_entities sp on et.type=sp.type
			   WHERE  (r.is_owner or  sp.type is not null) and (esm.owner_entity_id is not null)),

--getting already existing entity mappings
existing_entity_mappings as (select subs.id as subs_id,etmap.owner_entity_id as entity_id,etmap.id as entity_mapping_id,polgrp.name as pol_grp_name,polgrp.id as pol_grp_id from (select * from ac_pol_grp_subs) subs 
	inner join ac_pol_grps polgrp on polgrp.id=subs.ac_pol_grp_id and polgrp.entity_types @> array_append(array[]::integer[],@typeID) 
	inner join pol_grp_subs_entity_mappings etmap on etmap.pol_grp_subs_id=subs.id
	where subs.owner_space_id=@spaceID and subs.owner_team_id=@teamID)

	select et.*,etmap.policy_groups from (
		select et.entity_id,et.type,et.label from entities et left join owner_entities oe on oe.entity_id=et.entity_id left join subs_entities set on set.entity_id=et.entity_id %s
		group by et.entity_id %s
	) et
	left join lateral (
		select json_agg(json_build_object('pol_grp_name',polgrp.name,'pol_grp_id',polgrp.id,'entity_mapping_id',ext.entity_mapping_id)) as policy_groups from ac_pol_grps polgrp
		left join (select ext.* from existing_entity_mappings ext where ext.entity_id=et.entity_id) ext on ext.pol_grp_id=polgrp.id 
		where (
			ext.pol_grp_id is null and polgrp.is_predefined and polgrp.entity_types @> array_append(array[]::integer[],@typeID)
		) or (
			ext.pol_grp_id is not null
		) 
	)etmap on true`, policyGroupName, blocksFilterString, orderByString+limitString+offsetString)

	if err := db.Raw(query, valuesMap).Scan(&entitiesData.Entities).Error; err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Error getting entities for team", Response{}, http.StatusInternalServerError)
		return handlerResp
		// RespondWithError(w, http.StatusInternalServerError, "Error getting tags")
		// return
	}

	if err := db.Raw(countQuery, valuesMap).Scan(&entitiesData.Count).Error; err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Error getting entities for team", Response{}, http.StatusInternalServerError)
		return handlerResp
		// RespondWithError(w, http.StatusInternalServerError, "Error getting tags")
		// return
	}

	var resp Response
	resp.Data = entitiesData
	resp.Err = false
	resp.Msg = "Entities for the team retrieved successfully"

	handlerResp = common_services.BuildResponse(false, "entities for the team retrieved successfully", resp, http.StatusOK)

	return handlerResp
	// RespondWithJSON(w, http.StatusOK, resp)
}
