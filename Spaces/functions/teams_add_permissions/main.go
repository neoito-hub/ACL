package teams_add_permissions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aidarkhanov/nanoid"
	"github.com/neoito-hub/ACL-Block/Data-Models/models"
	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {

	var b RequestObject
	var handlerResp common_services.HandlerResponse
	var permissionIDs []string
	var existingPermissions []ExistingPermissions
	var polGrpSubs []map[string]interface{}
	var polGrpSubsEntityMappings []map[string]interface{}

	var deletedPolGrpSubsIDs []string
	var deletedEntityMappingIDs []string

	isOwner, isOwnerErr := strconv.ParseBool(payload.IsOwner)

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db
	valuesMap := make(map[string]interface{})
	permissionsMap := make(map[string]PermissionsMetadata)
	spaceAccessEntityMap := map[int]string{1: "1", 2: "2", 3: "3"}

	tx := db.Begin()

	for _, val := range b.Permissions {
		permissionIDs = append(permissionIDs, val.PermissionID)

	}

	availableEntitiesMap := make(map[string]interface{})
	allowedSpaceAccessEntitiesMap := make(map[int]bool)
	availableEntitiesMap["owner_user_id"] = payload.UserID
	availableEntitiesMap["owner_space_id"] = payload.SpaceID
	availableEntitiesMap["policy_group_list"] = []string{"ACL Access"}
	var allowedSpaceAccessEntities []int

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
	space_entities as (select distinct acl.type from acl_entities acl where acl.entity_id=acl.type::varchar)
	select * from space_entities `

	res := db.Raw(availableEntitiesQuery, availableEntitiesMap).Scan(&allowedSpaceAccessEntities)

	if res.Error != nil {

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	for _, val := range allowedSpaceAccessEntities {
		allowedSpaceAccessEntitiesMap[val] = true
	}

	valuesMap["owner_team_id"] = b.TeamID
	valuesMap["owner_space_id"] = payload.SpaceID
	valuesMap["permissionIDs"] = permissionIDs

	getExistingPermissionsQuery := `with etmap as 
	--getting entities as json_array
	(select subs.permission_id,json_agg(json_build_object('entity_id',et.entity_id,'entity_type',et.type,'entity_mapping_id',etmap.id)) as entities
	from (select subs.* from ac_pol_grp_subs subs where subs.owner_space_id=@owner_space_id and subs.owner_team_id=@owner_team_id
									)subs inner join pol_grp_subs_entity_mappings etmap on etmap.pol_grp_subs_id=subs.id 
	inner join entities et on et.entity_id=etmap.owner_entity_id group by subs.permission_id),
	existing_policy_groups as (
	select acper.id as permission_id,
	json_agg(json_build_object('subs_id',subs.id,'polgrp_id',polgrp.id,'entity_types',polgrp.entity_types,'new_subs_id',null)) as policy_groups from 
	ac_permissions acper inner join 
	 (select subs.* from ac_pol_grp_subs subs where subs.owner_space_id=@owner_space_id and subs.owner_team_id=@owner_team_id
									   )subs on acper.id=subs.permission_id
	left join ac_pol_grps polgrp on polgrp.id=subs.ac_pol_grp_id
	group by acper.id
	),
	predefined_policy_groups as (
		
		select acper.id as permission_id,json_agg(json_build_object('subs_id',null,'polgrp_id',polgrp.id,'entity_types',polgrp.entity_types,'new_subs_id',null)) 
as predefined_policy_groups  from 
	ac_permissions acper 
	left join per_pol_grps prpg on acper.id=prpg.ac_permission_id
	left join ac_pol_grps polgrp on prpg.ac_pol_grp_id=polgrp.id
	group by acper.id
	)
	select  acper.id,pg.policy_groups,etmap.entities,pre.predefined_policy_groups from ac_permissions acper left join existing_policy_groups pg on pg.permission_id=acper.id
	left join etmap on etmap.permission_id=acper.id left join predefined_policy_groups pre on pre.permission_id=acper.id
	`

	if err := tx.Raw(getExistingPermissionsQuery, valuesMap).Scan(&existingPermissions).Error; err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	for _, val := range existingPermissions {
		var permissionsMetadata PermissionsMetadata
		policyGroupsMap := make(map[string]ExistingPolicyGroups)
		entitiesMap := make(map[string]ExistingEntities)
		predefinedPolicyGroupsMap := make(map[string]ExistingPolicyGroups)
		var policyGroups []PolicyGroup
		var entities []Entities
		var predefinedPolicyGroups []PolicyGroup

		if val.PolicyGroups != nil {
			pgErr := json.Unmarshal(val.PolicyGroups, &policyGroups)

			if pgErr != nil {
				handlerResp = common_services.BuildErrorResponse(true, "Error getting existing permissions", Response{Msg: "Error getting existing permissions"}, http.StatusInternalServerError)
				return handlerResp
			}
			for _, val := range policyGroups {
				var polGrp ExistingPolicyGroups

				if existingPolGrp, polGrpExists := policyGroupsMap[val.PolGrpID]; polGrpExists {
					polGrp.EntityTypes = existingPolGrp.EntityTypes
					polGrp.PolGrpID = existingPolGrp.PolGrpID
					if len(val.SubsID) > 0 {
						polGrp.SubsIDs = append(existingPolGrp.SubsIDs, val.SubsID)
					}
				} else {
					polGrp.EntityTypes = val.EntityTypes
					polGrp.PolGrpID = val.PolGrpID
					if len(val.SubsID) > 0 {
						polGrp.SubsIDs = append(polGrp.SubsIDs, val.SubsID)
					}
				}

				policyGroupsMap[val.PolGrpID] = polGrp

			}
		}
		if val.Entities != nil {
			etErr := json.Unmarshal(val.Entities, &entities)

			if etErr != nil {
				handlerResp = common_services.BuildErrorResponse(true, "Error getting existing entities", Response{Msg: "Error getting existing entities"}, http.StatusInternalServerError)
				return handlerResp
			}

			for _, val := range entities {
				var entity ExistingEntities

				if existingEntity, entityExists := entitiesMap[val.EntityID]; entityExists {
					entity.ChangeType = existingEntity.ChangeType
					entity.EntityID = existingEntity.EntityID
					entity.EntityType = existingEntity.EntityType
					if len(val.EntityMappingID) > 0 {
						entity.EntityMappingIDs = append(existingEntity.EntityMappingIDs, val.EntityMappingID)
					}
				} else {
					entity.ChangeType = val.ChangeType
					entity.EntityID = val.EntityID
					entity.EntityType = val.EntityType
					if len(val.EntityMappingID) > 0 {
						entity.EntityMappingIDs = append(entity.EntityMappingIDs, val.EntityMappingID)
					}
				}

				entitiesMap[val.EntityID] = entity

			}

		}
		if val.PredefinedPolicyGroups != nil {
			pgErr := json.Unmarshal(val.PredefinedPolicyGroups, &predefinedPolicyGroups)

			if pgErr != nil {

				handlerResp = common_services.BuildErrorResponse(true, "Error getting predefined policy groups", Response{Msg: "Error getting predefined policy groups"}, http.StatusInternalServerError)
				return handlerResp
			}
			for _, val := range predefinedPolicyGroups {
				var polGrp ExistingPolicyGroups

				polGrp.EntityTypes = val.EntityTypes
				polGrp.PolGrpID = val.PolGrpID

				predefinedPolicyGroupsMap[val.PolGrpID] = polGrp

			}

		}
		permissionsMetadata.EntitiesMap = entitiesMap
		permissionsMetadata.PolicyGroupsMap = policyGroupsMap
		permissionsMetadata.PredefinedPolicyGroupsMap = predefinedPolicyGroupsMap
		permissionsMap[val.ID] = permissionsMetadata

	}

	for _, val := range b.Permissions {
		perMap, perExists := permissionsMap[val.PermissionID]

		if !perExists {
			handlerResp = common_services.BuildErrorResponse(true, "Invalid permission id", Response{Msg: "Invalid permission id"}, http.StatusInternalServerError)
			return handlerResp
		}

		//entities map for faster getting required entities for a given entity type for policy groups
		changedEntities := make(map[int][]Entities)

		for _, spaceEntity := range val.AddedSpaceAccessEntities {
			entityID, entityExists := spaceAccessEntityMap[spaceEntity.Type]

			_, allowedEntityAccess := allowedSpaceAccessEntitiesMap[spaceEntity.Type]

			if !allowedEntityAccess && (!isOwner || isOwnerErr != nil) {

				handlerResp = common_services.BuildErrorResponse(true, "No access permission", Response{}, http.StatusForbidden)
				return handlerResp
			}

			if entityExists {

				entity := Entity{ID: entityID, Type: spaceEntity.Type}

				val.AddedEntities = append(val.AddedEntities, entity)

			}
		}

		for _, spaceEntity := range val.DeletedSpaceAccessEntities {
			entityID, entityExists := spaceAccessEntityMap[spaceEntity.Type]

			if entityExists {
				entity := Entity{ID: entityID, Type: spaceEntity.Type}

				val.DeletedEntities = append(val.DeletedEntities, entity)

			}
		}

		//If permission is to be deleted as whole dont loop through added and deleted entities
		if val.IsDelete {
			for _, existingEntity := range perMap.EntitiesMap {
				deletedEntityMappingIDs = append(deletedEntityMappingIDs, existingEntity.EntityMappingIDs...)
			}

			for _, existingPolicyGroup := range perMap.PolicyGroupsMap {
				deletedPolGrpSubsIDs = append(deletedPolGrpSubsIDs, existingPolicyGroup.SubsIDs...)
			}
		} else {
			//For entities that are to be added to the permission
			for _, addedEntity := range val.AddedEntities {

				if _, entityExists := perMap.EntitiesMap[addedEntity.ID]; !entityExists {
					changedEntities[addedEntity.Type] = append(changedEntities[addedEntity.Type], Entities{EntityType: addedEntity.Type, EntityID: addedEntity.ID, ChangeType: 1, EntityMappingID: ""})
				}
			}

			//For entities that are to be deleted from the permission
			for _, deletedEntity := range val.DeletedEntities {
				//delete if only entity exists in the db
				if existingEntity, entityExists := perMap.EntitiesMap[deletedEntity.ID]; entityExists {
					deletedEntityMappingIDs = append(deletedEntityMappingIDs, existingEntity.EntityMappingIDs...)
				}
			}

		}

		//processing for policy groups to be added or deleted
		for _, polGrp := range perMap.PredefinedPolicyGroupsMap {

			var policyGroupSubIDs []string

			if existingPolicyGroup, polGrpExists := perMap.PolicyGroupsMap[polGrp.PolGrpID]; !polGrpExists {
				policyGroupSubIDs = append(policyGroupSubIDs, nanoid.New())

				newPolGrpSubs := map[string]interface{}{
					"ID":           policyGroupSubIDs[0],
					"PermissionID": val.PermissionID,
					"OwnerSpaceID": payload.SpaceID,
					"OwnerTeamID":  b.TeamID,
					"AcPolGrpID":   polGrp.PolGrpID,
					"CreatedAt":    time.Now(),
					"UpdatedAt":    time.Now(),
				}

				polGrpSubs = append(polGrpSubs, newPolGrpSubs)

			} else {
				policyGroupSubIDs = append(policyGroupSubIDs, existingPolicyGroup.SubsIDs...)
			}

			for _, subID := range policyGroupSubIDs {
				//looping through all the entity types of the given policy group
				for _, entityType := range polGrp.EntityTypes {
					newEntities, entityExists := changedEntities[int(entityType)]

					//checking if the entity exists for the given entity type
					if entityExists {
						//if entity exists loop all the entities for the given entity type
						for _, entity := range newEntities {
							//loop through all the sub ids for which entities are to be added

							newEntityMapping := map[string]interface{}{
								"ID":            nanoid.New(),
								"OwnerEntityID": entity.EntityID,
								"PolGrpSubsID":  subID,
								"CreatedAt":     time.Now(),
								"UpdatedAt":     time.Now(),
							}

							polGrpSubsEntityMappings = append(polGrpSubsEntityMappings, newEntityMapping)

						}
					}

				}

			}

		}
	}

	if len(polGrpSubs) > 0 {

		if err := tx.Model(&models.AcPolGrpSub{}).Create(&polGrpSubs).Error; err != nil {
			tx.Rollback()
			handlerResp = common_services.BuildErrorResponse(true, "Error creating policy group subscriptions", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

	}

	if len(polGrpSubsEntityMappings) > 0 {

		if err := tx.Model(&models.PolGrpSubsEntityMapping{}).Create(&polGrpSubsEntityMappings).Error; err != nil {
			tx.Rollback()
			handlerResp = common_services.BuildErrorResponse(true, "Error creating policy group subscription entity mappings", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

	}

	if len(deletedEntityMappingIDs) > 0 {

		if err := tx.Exec(`delete from pol_grp_subs_entity_mappings where id in ?`, deletedEntityMappingIDs).Error; err != nil {
			tx.Rollback()
			handlerResp = common_services.BuildErrorResponse(true, "Error deleting  policy group subscriptions", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

	}

	if len(deletedPolGrpSubsIDs) > 0 {

		if err := tx.Exec(`delete from ac_pol_grp_subs where id in ?`, deletedPolGrpSubsIDs).Error; err != nil {
			tx.Rollback()
			handlerResp = common_services.BuildErrorResponse(true, "Error deleting  policy group subscriptions", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

	}

	var resp []ResponseData

	tx.Commit()

	handlerResp = common_services.BuildResponse(false, "Permissions for team updated successfully!", Response{Data: resp, Err: false, Msg: "Permission subscription for team added successfully!"}, http.StatusOK)
	return handlerResp
}
