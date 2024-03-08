package shield_operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/neoito-hub/ACL-Block/captain/common_services"
	"gorm.io/gorm"
)

func VerifyAndGetUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) (common_services.ContextData, error) {
	// Validating and retreving user id from user access token
	bearToken := r.Header.Get("Authorization")

	client := &http.Client{}

	req, callerr := http.NewRequest("GET", os.Getenv("SHIELD_URL")+"/get-user-id", http.NoBody)
	if callerr != nil {
		return common_services.ContextData{}, callerr
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearToken)

	response, callerr := client.Do(req)

	if callerr != nil {
		// defer response.Body.Close()
		return common_services.ContextData{}, callerr
	}

	defer response.Body.Close()

	bodyBytes, callerr := io.ReadAll(response.Body)

	if callerr != nil {
		return common_services.ContextData{}, callerr
	}

	var responseObject ShiledResponse

	jsonErr := json.Unmarshal(bodyBytes, &responseObject)

	if jsonErr != nil {
		return common_services.ContextData{}, jsonErr
	}

	if !responseObject.Success {
		fmt.Println(responseObject.Data)
		return common_services.ContextData{}, errors.New(responseObject.Data.(string))
	}

	userID := (responseObject.Data.(map[string]interface{}))["user_id"].(string)

	var userData = common_services.ContextData{
		UserID: userID,
	}

	return userData, nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {

	var resp ErrResponse
	resp.Err = true
	resp.Msg = message
	RespondWithJSON(w, code, resp)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dataPayload, err := json.Marshal(payload)
	if err != nil {
		code = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("entered here")
	w.WriteHeader(code)

	w.Write(dataPayload)

}
