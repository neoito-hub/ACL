package roles_invite_user_email

import (
	"time"

	"gorm.io/datatypes"
)

type RequestObject struct {
	Data  []RolePayload `json:"data"`
	Email []string      `json:"email"`
}

type RolePayload struct {
	SpaceID string   `json:"space_id"`
	RoleIDs []string `json:"role_ids"`
}

type ExistingInviteEmails struct {
	Email string `json:"email"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Err bool   `json:"err"`
	Msg string `json:"msg"`
}

type ShiledResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

type ShieldUserData struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type EmailSendPayload struct {
	Link       string     `json:"link"`
	Email      string     `json:"email"`
	InviteInfo InviteData `json:"invite_data"`
}

type InviteData struct {
	SpaceData datatypes.JSON `json:"space_data"`
	RoleData  datatypes.JSON `json:"role_data"`
}

type LinkPayload struct {
	InviteID string `json:"invite_id"`
}

type EmailStruct struct {
	InviteLink  string
	SpaceName   string
	Email       string
	FirstLetter string
}

type InviteCreateResponse struct {
	InviteDetails      []CreatedInviteDetails `json:"invite_details"`
	ExistingRolesAdded []AddedRoleDetails     `json:"existing_roles_added"`
}

type AddedRoleDetails struct {
	OwnerUserID string `json:"user_id"`
	RoleID      string `json:"role_id"`
}

type CreatedInviteDetails struct {
	InvitedSpaceID string `json:"invited_space_id"`
	InvitedRoleID  string `json:"invited_role_id"`
}

type SpaceObject struct {
	SpaceID string `json:"space_id"`
	Exists  bool   `json:"exists"`
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
}

type ExistingObject struct {
	UserID string `json:"user_id"`
	Exists bool   `json:"exists"`
}

type RoleObject struct {
	RoleID string `json:"role_id"`
	Exists bool   `json:"exists"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type InviteObject struct {
	ID        string    `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
	Email     string    `json:"email"`
}
type SpaceData struct {
	SpaceID               string         `json:"space_id"`
	LegalID               string         `json:"legal_id"`
	Type                  string         `json:"type"`
	Name                  string         `json:"name"`
	BusinessName          string         `json:"business_name"`
	Address               string         `json:"address"`
	Email                 string         `json:"email"`
	Country               string         `json:"country"`
	BusinessCategory      string         `json:"business_category"`
	Description           string         `json:"description"`
	MetaData              datatypes.JSON `json:"metadata"`
	TaxPayerID            string         `json:"tax_payer_id"`
	DistinguishedName     string         `json:"distinguished_name"`
	Status                int            `json:"status"`
	MarketPlaceID         string         `json:"market_place_id"`
	DeveloperPortalAccess bool           `json:"developer_portal_access"`
}
