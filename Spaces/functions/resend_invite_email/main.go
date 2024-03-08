package resend_invite_email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
	"github.com/neoito-hub/ACL-Block/spaces/functions/mailer"
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

	var existingInvites []ExistingInvites

	invRes := db.Raw(`select i.id as invite_id, i.email, nanoid() as new_invite_id, s.space_id, s.name as space_name from invites i inner join invite_details id on id.invite_id=i.id inner join spaces s on s.space_id=id.invited_space_id left join (select u.email from users u inner join space_members m on m.owner_user_id=u.user_id where m.owner_space_id=?) u on u.email = i.email where invited_space_id=? and i.id in (?) and i.invite_type=1 and i.status != 2 and u.email is null group by i.id,s.space_id,s.name `, payload.SpaceID, payload.SpaceID, b.InviteIds).Scan(&existingInvites)

	if invRes.Error != nil {

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	if len(existingInvites) < 1 {
		handlerResp = common_services.BuildErrorResponse(true, "NO RECORD FOUND!", Response{}, http.StatusNoContent)
		return handlerResp
	}

	var inviteArray, invIds []string
	emailLinkMap := make(map[string]string)

	for _, v := range existingInvites {
		inviteArray = append(inviteArray, fmt.Sprintf("('%v','%v')", v.InviteID, v.NewInviteID))
		invIds = append(invIds, v.InviteID)

		linkErr, link := createInviteLink(LinkPayload{InviteID: v.NewInviteID})
		if linkErr != nil {
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

		emailLinkMap[v.Email] = link

	}

	var tx = db.Begin()

	if err := tx.Exec(fmt.Sprintf("insert into invites(id,email,invite_code,created_by,expires_at,created_at,updated_at,invite_type,status) select im.new_invite_id,i.email,i.invite_code,?,now() + interval'24hours',now(),now(),1,1 from invites i inner join (values%v) as im(invite_id, new_invite_id) on im.invite_id=i.id", strings.Join(inviteArray, ",")), userData.UserID).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	if err := tx.Exec(fmt.Sprintf("update invite_details id set invite_id = im.new_invite_id from (values%v) as im(invite_id, new_invite_id) where id.invite_id = im.invite_id", strings.Join(inviteArray, ","))).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	if err := tx.Exec("delete from invites where id in(?)", invIds).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	for email, link := range emailLinkMap {
		emailSendErr := SendUserInviteEmail(EmailSendPayload{Link: link, Email: email}, existingInvites[0].SpaceName)
		if emailSendErr != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

	}

	tx.Commit()
	var resp Response

	resp.Err = false
	resp.Msg = "Email sent successfully"

	if len(existingInvites) < len(b.InviteIds) {
		resp.Msg = "Email sent successfully. Existing members are excluded."
	}

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, resp.Msg, resp, http.StatusOK)
	return handlerResp
}

func createInviteLink(linkPayload LinkPayload) (error, string) {

	req, err := http.NewRequest("GET", os.Getenv("SPACE_URL")+"/invitation/", nil)
	if err != nil {
		return err, ""
	}

	q := req.URL.Query()

	q.Add("invite_id", linkPayload.InviteID)

	req.URL.RawQuery = q.Encode()

	return nil, req.URL.String()

	// req, err := http.NewRequest("POST", os.Getenv("SHIELD_URL")+"/login", nil)
	// if err != nil {
	// 	return err, ""
	// }

	// spaceLink := os.Getenv("USER_INVITE_ORG_URL") + fmt.Sprintf("?invite_id=%s", linkPayload.InviteID)

	// // if you appending to existing query this works fine
	// q := req.URL.Query()
	// q.Add("org_url", spaceLink)
	// q.Add("Client-Id", os.Getenv("SHIELD_CLIENT_ID"))

	// req.URL.RawQuery = q.Encode()

	// return nil, req.URL.String()

}

func SendUserInviteEmail(emailData EmailSendPayload, spaceName string) error {

	body, bodyerr := buildTemplate(EmailStruct{InviteLink: emailData.Link, SpaceName: spaceName, Email: emailData.Email, FirstLetter: strings.ToUpper(spaceName[0:1])})
	if bodyerr != nil {
		return bodyerr
	}

	_, emailerr := mailer.SendEmail(body, []string{emailData.Email})
	if emailerr != nil {
		log.Println("Email Not Send")

		return emailerr
	}

	return nil

}

// Function returns the email verification template with the provided user data
func buildTemplate(userMailData EmailStruct) (bytes.Buffer, error) {
	dir, direrr := os.Getwd()
	if direrr != nil {
		log.Fatal(direrr)
	}

	templateFilePath := dir + "/static/src/templates/user_invite.html"

	t, templateErr := template.ParseFiles(templateFilePath)
	if templateErr != nil {
		fmt.Printf("template error is %v", templateErr)
		return bytes.Buffer{}, templateErr
	}

	var body bytes.Buffer

	fromMail := fmt.Sprintf("From: Appblocks <%s>\r\n", os.Getenv("SHIELD_MAILER_EMAIL"))
	toMail := fmt.Sprintf("To: <%s>\r\n", userMailData.Email)

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subjectData := fmt.Sprintf("Subject: Appblocks Space Invite \n%s\n\n", mimeHeaders)

	mailData := fromMail + toMail + subjectData

	body.Write([]byte(mailData))

	fmt.Println("asfdasfd", t)

	terr := t.Execute(&body, userMailData)
	if terr != nil {
		fmt.Printf("error is %v", terr)
		return body, terr
	}

	return body, nil
}
