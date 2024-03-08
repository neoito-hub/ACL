package roles_delete_existing_pol_grp_subs

import (
	"encoding/json"
	"net/http"

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

	tx := db.Begin()
	subsValuesMap := make(map[string]interface{})
	subsValuesMap["subsID"] = b.ID

	existingMappingsDeleteQuery := `delete from pol_grp_subs_entity_mappings etmap where id in (select etmap.id as mapping_id from
		pol_grp_subs_entity_mappings etmap inner join ac_pol_grp_subs subs on subs.id=etmap.pol_grp_subs_id
		where subs.id = @subsID)`

	if err := tx.Exec(existingMappingsDeleteQuery, subsValuesMap).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp

	}

	res := tx.Where("id = ?", b.ID).Unscoped().Delete(&models.AcPolGrpSub{})
	if res.Error != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}
	tx.Commit()

	// RespondWithJSON(w, http.StatusOK, Response{
	// 	Err: false,
	// 	Msg: "Role policy group subscription removed successfully!",
	// })

	handlerResp = common_services.BuildResponse(false, "Role policy group subscription removed successfully!", Response{
		Err: false,
		Msg: "Role policy group subscription removed successfully!",
	}, http.StatusOK)
	return handlerResp
}
