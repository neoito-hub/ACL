package invite_user_email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/aidarkhanov/nanoid"
	"github.com/neoito-hub/ACL-Block/Data-Models/models"
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

	//temporary map to remove duplication
	spaceMap := make(map[string]map[string]bool)

	spaceTeamMap := make(map[string][]string)

	var inviteCodeArray []string
	// items for building payload for existing spaces fetch from db
	var spaceArray []string
	spaceCounter := 0
	spacesValuesMap := make(map[string]interface{})

	//items for building payload for existing teams fetch from db
	var teamArray []string
	teamCounter := 0
	teamsValuesMap := make(map[string]interface{})
	createEmailMap := make(map[string]bool)

	for _, space := range b.Data {
		if _, ok := spaceMap[space.SpaceID]; ok {

		} else {
			spaceMap[space.SpaceID] = make(map[string]bool)
			inviteCodeArray = append(inviteCodeArray, space.SpaceID)

		}
		for _, team := range space.TeamIDs {
			if _, ok := spaceMap[space.SpaceID][team]; ok {

			} else {
				spaceMap[space.SpaceID][team] = true
				spaceTeamMap[space.SpaceID] = append(spaceTeamMap[space.SpaceID], team)
			}
		}
	}

	sort.Strings(inviteCodeArray)

	inviteCode := ""

	for _, space := range inviteCodeArray {

		inviteCode += fmt.Sprintf("|%v|", space)

		sort.Strings(spaceTeamMap[space])

		teamsString := strings.Join(spaceTeamMap[space], "~")

		inviteCode += teamsString

	}

	for spaceID, space := range spaceMap {
		for _, email := range b.Email {
			createEmailMap[email] = true
			spaceCounter++
			one := fmt.Sprintf("space%v", spaceCounter)
			spaceCounter++
			two := fmt.Sprintf("space%v", spaceCounter)
			spacesValuesMap[one] = spaceID
			spacesValuesMap[two] = email

			spaceArray = append(spaceArray, fmt.Sprintf("(@%v,@%v)", one, two))

			for teamID, _ := range space {

				teamCounter++
				two := fmt.Sprintf("team%v", teamCounter)
				teamCounter++
				three := fmt.Sprintf("team%v", teamCounter)
				teamsValuesMap[two] = teamID
				teamsValuesMap[three] = email

				teamArray = append(teamArray, fmt.Sprintf("(@%v,@%v)", two, three))

			}
		}

	}

	//checking if same invites are already pending for the given emails

	var existingInviteEmails []ExistingInviteEmails

	db.Raw(`select i.email from invites i where i.status=1 and i.email in (?) and i.invite_type=1 and now()<i.expires_at and i.invite_code=?`, b.Email, inviteCode).Scan(&existingInviteEmails)

	if len(existingInviteEmails) > 0 {
		for _, existingEmail := range existingInviteEmails {
			fmt.Println(existingEmail.Email)
			delete(createEmailMap, existingEmail.Email)
		}

	}

	if len(createEmailMap) == 0 {
		var resp Response
		resp.Data = InviteCreateResponse{}

		resp.Err = false
		resp.Msg = "Invites Already exist for the given emails"

		// RespondWithJSON(w, http.StatusOK, resp)
		// return

		handlerResp = common_services.BuildResponse(false, "Invites Already exist for the given emails", resp, http.StatusOK)
		return handlerResp
	}

	var spaceDetails []SpaceObject
	var teamDetails []TeamObject
	var teamMember []models.TeamMember
	var inviteData []InviteObject
	var createdInviteDetails []CreatedInviteDetails
	var addedTeamDetails []AddedTeamDetails
	var inviteInfo InviteData

	// for building invite id array for getting metadata for sending emails
	var inviteIDArray []string

	//used for building payload for invite details table bulk create
	inviteDetailsCounter := 0
	var inviteDetailsArray []string
	inviteDetailsMap := make(map[string]interface{})

	//used for building payload for invite table bulk create
	inviteCounter := 0
	var inviteArray []string
	invitesMap := make(map[string]interface{})

	inviteIdEmailMap := make(map[string]string)

	var tx = db.Begin()

	var existingSpacesMap = map[string]map[string]ExistingObject{}
	var existingTeamsMap = map[string]map[string]ExistingObject{}

	spacesQuery := fmt.Sprintf("select v.space_id,true as exists,u.user_id as user_id,v.email from (values %v) as v(space_id,email) inner join space_members mr on mr.owner_space_id=v.space_id inner join users u on mr.owner_user_id=u.user_id and v.email=u.email", strings.Join(spaceArray, ","))

	teamsQuery := fmt.Sprintf("select v.team_id,true as exists,u.user_id as user_id,v.email from (values %v) as v(team_id,email) inner join team_members tm on tm.owner_team_id=v.team_id inner join users u on tm.member_id=u.user_id and v.email=u.email", strings.Join(teamArray, ","))

	// // TODO

	if err := tx.Raw(spacesQuery, spacesValuesMap).Scan(&spaceDetails).Error; err != nil {
		tx.Rollback()
		// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		// return

		handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	if len(teamArray) > 0 {
		if err := tx.Raw(teamsQuery, teamsValuesMap).Scan(&teamDetails).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}
	}

	for _, v := range spaceDetails {
		_, ok := existingSpacesMap[v.SpaceID]
		if !ok {
			existingSpacesMap[v.SpaceID] = map[string]ExistingObject{}
		}
		existingSpacesMap[v.SpaceID][v.Email] = ExistingObject{Exists: v.Exists, UserID: v.UserID}
	}

	for _, v := range teamDetails {
		_, ok := existingTeamsMap[v.TeamID]
		if !ok {
			existingTeamsMap[v.TeamID] = map[string]ExistingObject{}
		}
		existingTeamsMap[v.TeamID][v.Email] = ExistingObject{Exists: v.Exists, UserID: v.UserID}
	}

	//looping through the spaces provided in the payload
	for spaceID, teams := range spaceMap {

		for _, email := range b.Email {
			teamsCount := 0
			existingSpace, userExistsInSpace := existingSpacesMap[spaceID][email]

			//checking if there are teams for the given space id
			if len(teams) > 0 {

				//looping through the teams present for the above space
				for teamID, _ := range teams {
					_, userExistsInTeam := existingTeamsMap[teamID][email]
					if !userExistsInTeam {
						if userExistsInSpace {
							teamMember = append(teamMember, models.TeamMember{ID: nanoid.New(), OwnerTeamID: teamID, MemberID: existingSpace.UserID})
							teamsCount++
						} else {
							inviteDetailsCounter++
							one := fmt.Sprintf("invite%v", inviteDetailsCounter)
							inviteDetailsCounter++
							two := fmt.Sprintf("invite%v", inviteDetailsCounter)
							inviteDetailsCounter++
							three := fmt.Sprintf("invite%v", inviteDetailsCounter)
							inviteDetailsCounter++
							four := fmt.Sprintf("invite%v", inviteDetailsCounter)
							inviteDetailsCounter++
							five := fmt.Sprintf("invite%v", inviteDetailsCounter)

							inviteDetailsMap[one] = nanoid.New()
							inviteDetailsMap[two] = spaceID
							inviteDetailsMap[three] = email
							inviteDetailsMap[four] = teamID
							inviteDetailsMap[five] = createInvitePayload(&inviteCounter, email, &inviteArray, invitesMap, &inviteIDArray, inviteIdEmailMap, inviteCode, userData.UserID)

							inviteDetailsArray = append(inviteDetailsArray, fmt.Sprintf("(@%v,@%v,@%v,@%v,@%v,now(),now())", one, two, three, four, five))
							teamsCount++
						}
					}
				}
			}

			//if no teams are present
			if teamsCount == 0 {
				fmt.Printf("entered here")

				if !userExistsInSpace {
					inviteDetailsCounter++
					one := fmt.Sprintf("invite%v", inviteDetailsCounter)
					inviteDetailsCounter++
					two := fmt.Sprintf("invite%v", inviteDetailsCounter)
					inviteDetailsCounter++
					three := fmt.Sprintf("invite%v", inviteDetailsCounter)
					inviteDetailsCounter++
					four := fmt.Sprintf("invite%v", inviteDetailsCounter)

					inviteDetailsMap[one] = nanoid.New()
					inviteDetailsMap[two] = spaceID
					inviteDetailsMap[three] = email

					inviteDetailsMap[four] = createInvitePayload(&inviteCounter, email, &inviteArray, invitesMap, &inviteIDArray, inviteIdEmailMap, inviteCode, userData.UserID)

					inviteDetailsArray = append(inviteDetailsArray, fmt.Sprintf("(@%v,@%v,@%v,null,@%v,now(),now())", one, two, three, four))
				}
			}

		}
	}

	if len(teamMember) > 0 {
		if err := tx.Create(teamMember).Scan(&addedTeamDetails).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}
	}

	if len(inviteArray) > 0 {
		if err := tx.Raw(fmt.Sprintf("insert into invites(id,email,invite_code,created_by,expires_at,created_at,updated_at,invite_type,status) values %v returning *", strings.Join(inviteArray, ",")), invitesMap).Scan(&inviteData).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}
	}

	if len(inviteDetailsArray) > 0 {
		if err := tx.Raw(fmt.Sprintf("insert into invite_details(id,invited_space_id,email,invited_team_id,invite_id,created_at,updated_at) values %v returning *", strings.Join(inviteDetailsArray, ",")), inviteDetailsMap).Scan(&createdInviteDetails).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}
	}

	if len(inviteArray) > 0 {
		if err := tx.Raw(`select t.team_data,s.space_data from 
	--fetching team details
	(select json_agg(json_build_object('team_id',t.team_id,'team_name',t.name)) as team_data from invite_details i inner join teams t on t.team_id=i.invited_team_id where i.invite_id in (?))t 
	--fetching space data
	left join ( select json_agg(json_build_object('space_id',s.space_id,'space_name',s.name)) as space_data from invite_details i inner join spaces s on s.space_id=i.invited_space_id where i.invite_id in (?) )s on true `, inviteIDArray, inviteIDArray).Scan(&inviteInfo).Error; err != nil {
			tx.Rollback()
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}
	}

	emailLinkMap := make(map[string]string)

	for _, v := range inviteData {
		linkErr, link := createInviteLink(LinkPayload{InviteID: v.ID})
		if linkErr != nil {
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

		emailLinkMap[v.Email] = link

	}

	for email, link := range emailLinkMap {
		emailSendErr := SendUserInviteEmail(EmailSendPayload{Link: link, Email: email, InviteInfo: inviteInfo})
		if emailSendErr != nil {
			// RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			// return

			handlerResp = common_services.BuildErrorResponse(true, "Internal Server Error", Response{}, http.StatusInternalServerError)
			return handlerResp
		}

	}

	tx.Commit()
	var resp Response
	resp.Data = InviteCreateResponse{InviteDetails: createdInviteDetails, ExistingTeamsAdded: addedTeamDetails}

	resp.Err = false

	switch {
	case len(inviteArray) == 0:
		resp.Msg = "Member already exists."
	case len(inviteArray) < len(b.Email):
		resp.Msg = "Invites sent successfully. Existing members are excluded."
	default:
		resp.Msg = "Invites sent successfully."
	}

	// RespondWithJSON(w, http.StatusOK, resp)

	handlerResp = common_services.BuildResponse(false, resp.Msg, resp, http.StatusOK)
	return handlerResp
}

func createInvitePayload(inviteCounter *int, email string, inviteArray *[]string, invitesMap map[string]interface{}, inviteIDArray *[]string, inviteIdEmailMap map[string]string, inviteCode string, UserID string) string {

	var inviteID string

	if val, inviteExists := inviteIdEmailMap[email]; inviteExists {
		inviteID = val
	} else {
		inviteID = nanoid.New()
		inviteIdEmailMap[email] = inviteID
		*inviteCounter++
		one := fmt.Sprintf("invite%v", *inviteCounter)
		*inviteCounter++
		two := fmt.Sprintf("invite%v", *inviteCounter)
		*inviteCounter++
		three := fmt.Sprintf("invite%v", *inviteCounter)

		invitesMap[one] = inviteID
		invitesMap[two] = email
		invitesMap[three] = inviteCode

		*inviteArray = append(*inviteArray, fmt.Sprintf("(@%v,@%v,@%v,'%v',now() + interval'24hours',now(),now(),1,1)", one, two, three, UserID))

		*inviteIDArray = append(*inviteIDArray, inviteID)

	}

	return inviteID

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

func SendUserInviteEmail(emailData EmailSendPayload) error {
	var inviteData []interface{}

	err := json.Unmarshal(emailData.InviteInfo.SpaceData, &inviteData)

	if err != nil {
		return err
	}

	item, _ := inviteData[0].(map[string]interface{})
	fmt.Printf("item is %v", item)
	spaceName := item["space_name"].(string)
	fmt.Printf("space name is %v", spaceName)

	body, bodyerr := buildTemplate(EmailStruct{InviteLink: emailData.Link, SpaceName: spaceName, Email: emailData.Email, FirstLetter: strings.ToUpper(spaceName[0:1])})
	if bodyerr != nil {
		return bodyerr
	}

	_, emailerr := mailer.SendEmail(body, []string{emailData.Email})
	if emailerr != nil {
		log.Println("Email Not Send")
		log.Println(err)

		return err
	}

	return err

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
