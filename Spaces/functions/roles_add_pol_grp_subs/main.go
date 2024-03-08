package roles_add_pol_grp_subs

import (
	"encoding/json"
	"net/http"

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

	valuesMap := make(map[string]interface{})
	idMap := make(map[string]string)

	for _, id := range b.AcPolGrpIDs {
		idMap[id] = nanoid.New()
	}

	// db := DBInit()
	// sqlDB, dberr := db.DB()

	// if dberr != nil {
	// 	log.Fatalln(dberr)
	// }
	// defer sqlDB.Close()

	tx := db.Begin()

	var insertData []models.AcPolGrpSub

	valuesMap["role_id"] = b.RoleID
	valuesMap["owner_space_id"] = b.SpaceID
	valuesMap["ac_pol_grp_ids"] = b.AcPolGrpIDs

	selectQuery := `select @owner_space_id as owner_space_id,@role_id as role_id,pg.id as ac_pol_grp_id from ac_pol_grps pg LEFT JOIN (select * from ac_pol_grp_subs where permission_id is null) pgs ON pgs.ac_pol_grp_id=pg.id AND pgs.owner_space_id = @owner_space_id AND pgs.role_id = @role_id WHERE pg.id IN @ac_pol_grp_ids AND pgs.id is null`

	if err := tx.Raw(selectQuery, valuesMap).Scan(&insertData).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	var resp []ResponseData

	for i, p := range insertData {
		insertData[i].ID = nanoid.New()
		resp = append(resp, ResponseData{ID: insertData[i].ID, SpaceID: p.OwnerSpaceID, RoleID: p.RoleID, AcPolGrpID: p.AcPolGrpID})
	}

	// subsCreationQuery := `INSERT INTO ac_pol_grp_subs(
	// 	id, created_at, updated_at, deleted_at, owner_space_id, role_id, owner_team_id, ac_pol_grp_id, opt_counter, owner_user_id) VAL returning *;`

	// var ResponseData []models.AcPolGrpSub

	if err := tx.Create(&insertData).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	// if err := tx.Create(insertData,  valuesMap).Scan(&ResponseData).Error; err != nil {
	// 	tx.Rollback()
	// 	RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	// 	return
	// }

	// res := tx.Where("id = ?", b.ID).Unscoped().Delete(&models.AcPolGrpSub{})
	// if res.Error != nil {
	// 	tx.Rollback()
	// 	RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	// 	return
	// }
	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{Data: resp, Err: false, Msg: "Policy subscription for roles added successfully!"})

	handlerResp = common_services.BuildResponse(false, "Policy subscription for roles updated successfully!", Response{Data: resp, Err: false, Msg: "Policy subscription for roles added successfully!"}, http.StatusOK)
	return handlerResp
}
