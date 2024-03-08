package teams_add_entities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aidarkhanov/nanoid"
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
	// _, shieldVerifyError := VerifyAndGetUser(w, r)
	// if shieldVerifyError != nil {
	// 	fmt.Printf("shieldVerifyError: %v\n", shieldVerifyError)
	// 	RespondWithError(w, http.StatusUnauthorized, shieldVerifyError.Error())

	// 	return
	// }

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db

	// db := DBInit()
	// sqlDB, dberr := db.DB()

	// if dberr != nil {
	// 	log.Fatalln(dberr)
	// }
	// defer sqlDB.Close()

	tx := db.Begin()

	var existingMappings ExistingMappings

	var polGrpInsertData []models.AcPolGrpSub
	var entityMappingInsertData []models.PolGrpSubsEntityMapping
	newEntityMap := make(map[string]NewEntityMappings)
	newPolicyGrpSubsMap := make(map[string]string)

	existingValuesMap := make(map[string]interface{})

	existingValuesMap["teamID"] = b.TeamID
	existingValuesMap["spaceID"] = payload.SpaceID

	existingEntitiesQuery := `
	select subs.polgrpsubs,etmap.etmappings from 
		(
			select json_object_agg(subs.ac_pol_grp_id,subs.id) as polgrpsubs from ac_pol_grp_subs subs
			where subs.owner_team_id=@teamID and subs.owner_space_id=@spaceID
		 ) subs
		 left join (
			 select json_object_agg(etmap.owner_entity_id,etmap.polgrpids) as etmappings from
			 (
				 select etmap.owner_entity_id,string_agg(subs.ac_pol_grp_id,',') as polgrpids from pol_grp_subs_entity_mappings etmap
				 inner join ac_pol_grp_subs subs on subs.id=etmap.pol_grp_subs_id
				 where subs.owner_team_id=@teamID and subs.owner_space_id=@spaceID
				 group by etmap.owner_entity_id
			 )etmap
		 )etmap on true
	`

	if err := tx.Raw(existingEntitiesQuery, existingValuesMap).Scan(&existingMappings).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp

	}

	grps := make(map[string]string)
	entityMappings := make(map[string]string)

	if existingMappings.Polgrpsubs != nil {
		err := json.Unmarshal(existingMappings.Polgrpsubs, &grps)

		if err != nil {
			tx.Rollback()
			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

		}
	}

	if existingMappings.Etmappings != nil {
		err := json.Unmarshal(existingMappings.Etmappings, &entityMappings)

		if err != nil {
			tx.Rollback()
			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

		}
	}

	// removing duplicates from new entity mappings array
	for _, val := range b.NewEntityMappings {

		// checking if the entity mapping already exists
		polgrpIDs, entityExists := entityMappings[val.EntityID]
		polGrpExists := strings.Contains(polgrpIDs, val.AcPolGrpID)
		_, entityMappingExists := newEntityMap[fmt.Sprintf("%s/%s", val.AcPolGrpID, val.EntityID)]

		if !entityMappingExists && (!entityExists || !polGrpExists) {
			newEntityMap[fmt.Sprintf("%s/%s", val.AcPolGrpID, val.EntityID)] = val
		}

		existingSubsID, subExists := grps[val.AcPolGrpID]
		_, subsMapExists := newPolicyGrpSubsMap[fmt.Sprintf("%s", val.AcPolGrpID)]

		if !subsMapExists {
			if subExists {
				newPolicyGrpSubsMap[fmt.Sprintf("%s", val.AcPolGrpID)] = existingSubsID
			} else {
				subsID := nanoid.New()
				newPolicyGrpSubsMap[fmt.Sprintf("%s", val.AcPolGrpID)] = subsID
				polGrpInsertData = append(polGrpInsertData, models.AcPolGrpSub{
					ID:           subsID,
					OwnerTeamID:  b.TeamID,
					AcPolGrpID:   val.AcPolGrpID,
					OwnerSpaceID: payload.SpaceID,
				})
			}
		}

	}

	for _, etmapping := range newEntityMap {
		polGrpSubsID, _ := newPolicyGrpSubsMap[fmt.Sprintf("%s", etmapping.AcPolGrpID)]

		entityMappingInsertData = append(entityMappingInsertData, models.PolGrpSubsEntityMapping{ID: nanoid.New(), OwnerEntityID: etmapping.EntityID, PolGrpSubsID: polGrpSubsID})

	}

	mappingValuesMap := make(map[string]interface{})
	subsValuesMap := make(map[string]interface{})

	mappingValuesMap["deletedIDs"] = b.DeletedEntityMappings
	subsValuesMap["teamID"] = b.TeamID
	mappingValuesMap["teamID"] = b.TeamID
	mappingValuesMap["spaceID"] = payload.SpaceID
	subsValuesMap["spaceID"] = payload.SpaceID
	subsValuesMap["typeID"] = b.TypeID

	var toDeleteSubsIds []string

	existingMappingsDeleteQuery := `delete from pol_grp_subs_entity_mappings etmap using (select etmap.id as mapping_id from
		pol_grp_subs_entity_mappings etmap inner join ac_pol_grp_subs subs on subs.id=etmap.pol_grp_subs_id
		where etmap.id in @deletedIDs 
		and subs.owner_team_id=@teamID and subs.owner_space_id=@spaceID
		)subs where subs.mapping_id=etmap.id returning etmap.pol_grp_subs_id`

	emptyPolicyGroupsDeleteQuery := `delete from ac_pol_grp_subs subs 
		using (select subs.id,count(etmap.id) as mapping_count from ac_pol_grp_subs subs inner join 
			  ac_pol_grps pol_grp on pol_grp.id=subs.ac_pol_grp_id inner join
			  pol_grp_subs_entity_mappings etmap on etmap.pol_grp_subs_id=subs.id
			  and subs.owner_user_id=@teamID and pol_grp.entity_type=1 and subs.owner_space_id=@spaceID group by subs.id) selsubs
		where subs.id=selsubs.id and selsubs.mapping_count=0 and subs.id in @toDeleteSubsIds`

	if len(b.DeletedEntityMappings) > 0 {
		if err := tx.Raw(existingMappingsDeleteQuery, mappingValuesMap).Scan(&toDeleteSubsIds).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp

		}

		subsValuesMap["toDeleteSubsIds"] = toDeleteSubsIds

		if err := tx.Exec(emptyPolicyGroupsDeleteQuery, subsValuesMap).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}
	}

	if len(newEntityMap) > 0 {
		if len(polGrpInsertData) > 0 {
			if err := tx.Create(&polGrpInsertData).Error; err != nil {
				tx.Rollback()
				// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
				// return

				handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
				return handlerResp
			}
		}

		if err := tx.Create(&entityMappingInsertData).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}
	}
	// tx.Rollback()

	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{Data: resp, Err: false, Msg: "Policy subscription for user added successfully!"})

	handlerResp = common_services.BuildResponse(false, "Entities for team added successfully", Response{Data: Response{}, Err: false, Msg: "Entities for team added successfully"}, http.StatusOK)
	return handlerResp
}
