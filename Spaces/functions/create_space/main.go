package create_space

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/aidarkhanov/nanoid"
	"github.com/neoito-hub/ACL-Block/Data-Models/models"
	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {

	// // Validating and retreving user id from user access token
	// userData, shieldVerifyError := VerifyAndGetUser(w, r)
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

	// Validating request body and method
	if len(b.Type) == 0 || (b.Type == "B" && (b.Country == "" || b.BusinessName == "" || b.Address == "")) {
		handlerResp = common_services.BuildErrorResponse(true, "Missing values", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")

	if !re.MatchString(b.Name) || !re.MatchString(b.BusinessName) {
		handlerResp = common_services.BuildErrorResponse(true, "No special characters allowed other than underscore for name and business name", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	userData := ShieldUserData{
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}

	tx := db.Begin()

	// space and member common ID
	SpaceID := nanoid.New()
	// role ID
	RoleID := nanoid.New()

	// Create member
	member := models.Member{
		ID:   SpaceID,
		Type: "S",
	}

	memberResult := tx.Create(&member)

	if memberResult.Error != nil {
		log.Println("!!! Error creating member => ", memberResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating member")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating member", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	// Create Space
	space := models.Space{
		SpaceID:               SpaceID,
		Type:                  b.Type,
		Name:                  b.Name,
		Email:                 b.Email,
		Country:               b.Country,
		BusinessCategory:      b.BusinessCategory,
		Description:           b.Description,
		MarketPlaceID:         b.MarketPlaceID,
		DeveloperPortalAccess: b.DeveloperPortalAccess,
		BusinessName:          b.BusinessName,
		Address:               b.Address,
	}

	spaceResult := tx.Create(&space)

	if spaceResult.Error != nil {
		log.Println("!!! Error creating space => ", spaceResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating space")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating space", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	// create space member
	spaceMember := models.SpaceMember{
		ID:           nanoid.New(),
		OwnerUserID:  userData.UserID,
		OwnerSpaceID: SpaceID,
	}

	spaceMemberResult := tx.Create(&spaceMember)

	if spaceMemberResult.Error != nil {
		log.Println("!!! Error creating space member => ", spaceMemberResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating space member")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating space member", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	// Create role
	roles := models.Role{
		ID:           RoleID,
		IsOwner:      true,
		Name:         "owner",
		OwnerSpaceID: SpaceID,
		CreatedBy:    userData.UserID,
		UpdatedBy:    userData.UserID,
	}

	roleResult := tx.Create(&roles)

	if roleResult.Error != nil {
		log.Println("!!! Error creating role => ", roleResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating role")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating role", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	// Create member role
	memberRole := models.MemberRole{
		ID:           nanoid.New(),
		RoleID:       RoleID,
		OwnerUserID:  userData.UserID,
		OwnerSpaceID: SpaceID,
	}

	memberRoleResult := tx.Create(&memberRole)

	if memberRoleResult.Error != nil {
		log.Println("!!! Error creating Member => ", memberRoleResult.Error)
		tx.Rollback()
		// RespondWithError(w, http.StatusBadRequest, "Error creating member role")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating member role", Response{}, http.StatusBadRequest)
		return handlerResp
	}

	//deprecated space manage purchase access
	// var Exists subExists

	// valuesMap := make(map[string]interface{})
	// valuesMap["userID"] = userData.UserID
	// valuesMap["spaceID"] = SpaceID

	// existsResult := tx.Raw(`select exists(select subs.id from ac_pol_grp_subs subs
	// 	inner join ac_pol_grps agp on agp.id=subs.ac_pol_grp_id
	// 	where  agp.name in ('Blocks-Manage-Purchase','Blocks Publish') and agp.type=1
	// 	 and subs.owner_space_id=@spaceID and subs.owner_user_id=@userID
	// 	) as pol_grp_subs_exists`, valuesMap).Scan(&Exists)

	// if existsResult.Error != nil {
	// 	log.Println("!!! Error on fetching exist check => ", existsResult.Error)
	// 	tx.Rollback()

	// 	handlerResp = common_services.BuildErrorResponse(true, "Error on fetching exist check", Response{}, http.StatusBadRequest)
	// 	return handlerResp
	// }

	// if !Exists.PolGrpSubsExists {
	// 	result := tx.Exec(`INSERT INTO public.ac_pol_grp_subs(
	// 		id, created_at, updated_at, owner_space_id, role_id, owner_team_id, owner_user_id, ac_pol_grp_id)
	// 		select nanoid(),now(),now(),@spaceID,null,null,@userID,polgrp.id
	// 		from ac_pol_grps polgrp where
	// 		polgrp.name in ('Blocks-Manage-Purchase','Blocks Publish') and polgrp.is_predefined and polgrp.type=1`, valuesMap)

	// 	if result.Error != nil {
	// 		log.Println("!!! Error creating ac_pol_grp_subs => ", result.Error)
	// 		tx.Rollback()

	// 		handlerResp = common_services.BuildErrorResponse(true, "Error creating ac_pol_grp_subs", Response{}, http.StatusBadRequest)
	// 		return handlerResp
	// 	}

	// }

	// setPolErr := set_accesspolicy_for_role.AssignPredefinedPolicyToRole(tx, set_accesspolicy_for_role.RolePayload{
	// 	SpaceID: SpaceID,
	// 	RoleID:  RoleID,
	// })

	// if setPolErr != nil {
	// 	log.Println("!!! Error setting polices => ", setPolErr)
	// 	RespondWithError(w, http.StatusBadRequest, "Error setting polices")
	// 	tx.Rollback()

	// 	return
	// }

	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{Err: false, Msg: "Successfully created space"})

	handlerResp = common_services.BuildResponse(false, "Successfully created space", Response{Err: false, Msg: "Successfully created space"}, http.StatusOK)
	return handlerResp
}
