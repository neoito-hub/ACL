package create_entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func ValidateRequest(w http.ResponseWriter, r *http.Request) (RequestObject, error) {
	var b RequestObject

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&b); err != nil {
		return b, err
	}

	if r.Method == http.MethodOptions {
		return b, errors.New("invalid method")
	}

	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(b.Name)

	if !isAlphanumeric {
		return b, errors.New("name should be alphanumeric")
	}

	return b, nil
}

// func structLoop(b RequestObject) {
// 	v := reflect.ValueOf(b)
// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)
// 		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
// 			fmt.Println(field)
// 			// bar
// 			// 1
// 		}
// 	}
// }

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
