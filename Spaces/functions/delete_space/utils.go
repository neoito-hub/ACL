package delete_space

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
)

func ValidateRequest(w http.ResponseWriter, r *http.Request) (RequestObject, error) {
	var b RequestObject

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&b); err != nil {
		fmt.Println(err)
		return b, err
	}

	if r.Method == http.MethodOptions {
		return b, errors.New("invalid method")
	}

	return b, nil
}

func VerifyAndGetUser(w http.ResponseWriter, r *http.Request) (ShieldUserData, error) {
	// Validating and retreving user id from user access token
	var userData ShieldUserData

	bearToken := r.Header.Get("Authorization")

	client := &http.Client{}

	req, callerr := http.NewRequest("GET", os.Getenv("SHIELD_URL")+"/get-user-id", http.NoBody)
	if callerr != nil {
		return userData, callerr
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearToken)

	response, callerr := client.Do(req)

	if callerr != nil {
		// defer response.Body.Close()
		return userData, callerr
	}

	defer response.Body.Close()

	bodyBytes, callerr := io.ReadAll(response.Body)

	if callerr != nil {
		return userData, callerr
	}

	var responseObject ShiledResponse

	jsonErr := json.Unmarshal(bodyBytes, &responseObject)

	if jsonErr != nil {
		return userData, jsonErr
	}

	if !responseObject.Success {
		fmt.Println(responseObject.Data)
		return userData, errors.New(responseObject.Data.(string))
	}

	userID, _ := (responseObject.Data.(map[string]interface{}))["user_id"].(string)

	userData = ShieldUserData{
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
	w.WriteHeader(code)

	w.Write(dataPayload)

}

func DeleteQueryRun(tx *gorm.DB, query string, value interface{}) error {
	var deleteResult interface{}

	deleteData := tx.Raw(query, value).Scan(&deleteResult)
	if deleteData.Error != nil {
		fmt.Println("=> Delete Error ============ ", deleteData.Error)
		return deleteData.Error
	} else if deleteData.RowsAffected < 1 {
		fmt.Println("=> No data to delete for query")
		return nil
	}

	return nil
}
