package roles_create_invite_link

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

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

	userData := ShieldUserData{
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}

	var spaceArray []string

	spaceRoleMap := make(map[string]map[string]bool)
	spaceMap := make(map[string][]string)

	for _, space := range b.Data {
		if _, ok := spaceRoleMap[space.SpaceID]; ok {

		} else {
			spaceRoleMap[space.SpaceID] = make(map[string]bool)
			spaceArray = append(spaceArray, space.SpaceID)

		}

		for _, role := range space.RoleIDs {

			if _, ok := spaceRoleMap[space.SpaceID][role]; ok {

			} else {
				spaceRoleMap[space.SpaceID][role] = true
				spaceMap[space.SpaceID] = append(spaceMap[space.SpaceID], role)
			}

		}

	}

	sort.Strings(spaceArray)

	inviteCode := ""

	for _, space := range spaceArray {

		inviteCode += fmt.Sprintf("|%v|", space)

		sort.Strings(spaceMap[space])

		rolesString := strings.Join(spaceMap[space], "~")

		inviteCode += rolesString

	}

	// check user is owner of all spaces

	// var ownerSpaceIds []SpaceIds

	// res := db.Raw("select DISTINCT(s.space_id) FROM spaces s INNER JOIN member_roles mr ON mr.owner_space_id = s.space_id INNER JOIN roles r ON r.owner_space_id = s.space_id AND r.id = mr.role_id WHERE s.space_id IN ? AND mr.owner_user_id = ? AND r.is_owner = true", spaceArray, shieldUser.UserID).Scan(&ownerSpaceIds)

	// if res.Error != nil {
	// 	RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	// 	return
	// }

	// if len(ownerSpaceIds) < len(spaceArray) {
	// 	RespondWithError(w, http.StatusUnauthorized, "You don't have the permissioin to create link")
	// 	return
	// }

	// check for link already exists for this combinatiion

	var linkResp LinkResponse
	res := db.Raw("select invite_link FROM invites WHERE invite_type = 2 AND invite_code = ?", inviteCode).Scan(&linkResp)
	if res.Error != nil {
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	if res.RowsAffected > 0 {
		db.Model(&models.Invites{}).Where("invite_type = 2 AND invite_code = ?", inviteCode).Update("expires_at", time.Now().Add(time.Hour*24))
		var resp Response
		resp.Data = linkResp
		resp.Err = false
		resp.Msg = "Link fetched successfully!"

		// RespondWithJSON(w, http.StatusOK, resp)
		// return

		handlerResp = common_services.BuildResponse(false, "Link fetched successfully!", resp, http.StatusOK)
		return handlerResp
	}

	inviteId := nanoid.New()

	linkErr, link := createInviteLink(LinkPayload{InviteID: inviteId})
	if linkErr != nil {
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	var inviteData LinkPayload

	var tx = db.Begin()

	// creating new entry in invites

	if err := tx.Raw("insert into invites(id,created_by,expires_at,created_at,updated_at,invite_type,invite_link,invite_code) values(?,?,now() + interval'24hours',now(),now(),2,?,?) returning *", inviteId, userData.UserID, link, inviteCode).Scan(&inviteData).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	// creating new entries in invite_details

	var invite_details []models.InviteDetails

	for _, space := range spaceArray {
		var invite_d models.InviteDetails
		roleExisted := false
		for _, role := range spaceMap[space] {

			invite_d.ID = nanoid.New()
			invite_d.InviteID = inviteId
			invite_d.InvitedSpaceID = space
			invite_d.InvitedRoleID = role

			roleExisted = true
			invite_details = append(invite_details, invite_d)

		}

		if !roleExisted {
			invite_d.ID = nanoid.New()
			invite_d.InviteID = inviteId
			invite_d.InvitedSpaceID = space

			invite_details = append(invite_details, invite_d)
		}

	}

	if err := tx.Create(&invite_details).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	tx.Commit()

	var resp Response
	resp.Data = LinkResponse{InviteLink: link}
	resp.Err = false
	resp.Msg = "Link generated successfully!"

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, "Link generated successfully!", resp, http.StatusOK)
	return handlerResp
}

func createInviteLink(linkPayload LinkPayload) (error, string) {
	req, err := http.NewRequest("GET", os.Getenv("SPACE_URL")+"/invitation/", nil)
	if err != nil {
		return err, ""
	}

	// spaceLink := os.Getenv("USER_INVITE_ORG_URL") + fmt.Sprintf("?invite_id=%s", linkPayload.InviteID)

	// if you appending to existing query this works fine
	q := req.URL.Query()
	// q.Add("org_url", spaceLink)
	// q.Add("client_id", os.Getenv("SHIELD_CLIENT_ID"))
	// q.Add("response_type", "code")
	// q.Add("state", "")

	q.Add("invite_id", linkPayload.InviteID)

	req.URL.RawQuery = q.Encode()

	return nil, req.URL.String()

}
