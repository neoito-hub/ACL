package gateway

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/neoito-hub/ACL-Block/captain/common_services"
	"github.com/neoito-hub/ACL-Block/captain/services/interface_layer"
	"github.com/neoito-hub/ACL-Block/captain/services/shield_operations"
	"gorm.io/gorm"
)

func Call(w http.ResponseWriter, r *http.Request, db *gorm.DB, resourcesMap map[string]common_services.Resources, actionName string, functionName string, spaceID string) (error, common_services.ContextData) {
	var shieldUser = common_services.ContextData{UserID: "", UserName: ""}
	var shieldVerifyError error
	resource, resourceExists := resourcesMap[functionName]
	bearToken := r.Header.Get("Authorization")

	// resource.IsAuthorised = 1
	// resource.IsAuthenticated = 1
	// resourceExists = true
	// shieldUser.UserID = "GJEF3Y6zlkLX4BnHx2Wgw"
	// shieldUser.UserID = "GYba44Jb1b9rJdWA4sq22"

	if !resourceExists {
		shield_operations.RespondWithError(w, http.StatusUnauthorized, "authentication error")
		return errors.New("authentication error"), common_services.ContextData{}
	}

	// authentication with shield in different ways depending on the type of api
	// for value of 3 may or maynot be authenticated depending on whether token is provided or not.
	// for any value other than 1 will be authenticated
	switch resource.IsAuthenticated {
	case 1:
	case 3:
		if len(bearToken) > 0 {
			shieldUser, shieldVerifyError = shield_operations.VerifyAndGetUser(w, r, db)
			if shieldVerifyError != nil {
				fmt.Printf("shieldVerifyError: %v\n", shieldVerifyError)
				shield_operations.RespondWithError(w, http.StatusUnauthorized, shieldVerifyError.Error())
				return errors.New("authentication error"), common_services.ContextData{}
			}
		}
	default:
		shieldUser, shieldVerifyError = shield_operations.VerifyAndGetUser(w, r, db)
		if shieldVerifyError != nil {
			fmt.Printf("shieldVerifyError: %v\n", shieldVerifyError)
			shield_operations.RespondWithError(w, http.StatusUnauthorized, shieldVerifyError.Error())

			return errors.New("authentication error"), common_services.ContextData{}
		}
	}

	if resource.IsAuthorised != 1 {

		if len(bearToken) > 0 {

			//authorisation of the resource using the interface layer
			ownerCheckData, ilErr := interface_layer.Handler(shieldUser.UserID,
				r, w, db, actionName, functionName, spaceID)

			if ilErr != nil {
				fmt.Printf("Interface layer Error: %v\n", ilErr.Error())
				interface_layer.RespondWithError(w, http.StatusForbidden, ilErr.Error())

				return errors.New("authorisation error"), common_services.ContextData{}
			}
			shieldUser.IsOwner = ownerCheckData.Exists
		}
	}
	shieldUser.SpaceID = spaceID

	return nil, shieldUser

}

func RespondWithJSON(w http.ResponseWriter, code int, payload string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}
