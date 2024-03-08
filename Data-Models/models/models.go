// UNION OF BLOCK REGISTRY (feat/block-discovery) model WITH PAYMENTS (develop) model (Jan 17 2023)

package models

import (
	"time"

	"github.com/lib/pq"

	// "github.com/shopspring/decimal"
	"gorm.io/datatypes"

	"gorm.io/gorm"
)

type Member struct {
	gorm.Model

	ID         string `gorm:"primaryKey;not null"`
	Type       string `gorm:"size:3"` //  (S-space, U-user, T-teams)
	OptCounter int    `gorm:"size:8"`
}

// Model for User Table
type User struct {
	UserID                  string `gorm:"primaryKey; size:255"` // FK to MEMBER table (type U).
	UserName                string `gorm:"size:255; not null; unique"`
	FullName                string `gorm:"size:255"`
	Email                   string `gorm:"size:255; not null; unique"`
	Password                string `gorm:"size:255; not null"`
	Address1                string `gorm:"size:150"`
	Address2                string `gorm:"size:150"`
	Phone                   string `gorm:"size:20"`
	EmailVerificationCode   string `gorm:"size:6"`
	EmailVerified           bool   `gorm:"default:false"`
	EmailVerificationExpiry time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
	OptCounter              int `gorm:"size:8"`

	UserMember Member `gorm:"foreignKey:UserID;References:ID"`
}

// shield models

type ShieldApp struct {
	//ID           uint
	AppId        string         `gorm:"primary_key; size:255"`
	ClientId     string         `gorm:"size:60; not null; unique"`
	ClientSecret string         `gorm:"size:255; not null"`
	UserId       string         `gorm:"size:255"`
	AppName      string         `gorm:"size:100; not null; unique"`
	AppSname     string         `gorm:"size:50"`
	Description  string         `gorm:"size:255"`
	LogoUrl      string         `gorm:"size:255"`
	AppUrl       string         `gorm:"size:255; not null"`
	RedirectUrl  pq.StringArray `gorm:"type:text[]; size:255; not null"`
	AppType      int            `gorm:"default:4"` // 1 - appblocks, 2 - internal app, 3 - client app, 4 - appblocks app
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	OwnerSpaceID string
	ID           int

	OwnerSpace Space `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
}

type Space struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SpaceID               string         `gorm:"primaryKey;not null"` // NOT NULL	FK to MEMBER table for this space or space unit (type S).
	LegalID               string         // The registered space identifier, given to the space (such as assigned by the government). This may be null for an space unit. This is not the name of the space, which should be stored in the ORGENTITYNAME table.
	Type                  string         `gorm:"not null"`             // NOT NULL P -> personal, B -> business or institution
	Name                  string         `gorm:"unique; not null"`     // Name of the space
	BusinessName          string         `gorm:"unique; default:null"` // The business name of the space. (Unique if value exist)
	Address               string         // Address of the business or institution
	LogoURL               string         // Space logo url
	Email                 string         // Email of the space (Unique if value exist)
	Country               string         // Country of the space (Not null for business or institution)
	BusinessCategory      string         // The business category, which describes the kind of business performed by an Space.
	Description           string         // A description of the Space.
	MetaData              datatypes.JSON // Additional metadata about the Space.
	TaxPayerID            string         // A string used to identify the Space for taxation purpose. Addition of this column triggered by Taxware integration, but presumably this column is useful even outside of Taxware.
	DistinguishedName     string         // Distinguished name (DN) of the Space. If LDAP is used, contains the DN of the Space in the LDAP server. If database is used as member repository, contains a unique name as defined by the membership hierarchy. DNs for all OrgEntities are logically unique, however due to the large field size, this constraint is not enforced at the physical level. The DN should not contain any spaces before or after the comma (,) or equals sign (=) and must be entered in lowercase.
	Status                int            // NOT NULL DEFAULT 0	The STATUS column in Space table indicates whether or not the space is locked. Valid values are as follows: 0 = not locked -1 = locked
	OptCount              int            //	The optimistic concurrency control counter for the table. Every time there is an update to the table, the counter is incremented.
	MarketPlaceID         string         // The ID of the market place where the space is registered.
	DeveloperPortalAccess bool           `gorm:"not null;default:false"` // Indicates whether the space has access to the developer portal.
	OptCounter            int            `gorm:"size:8"`

	Member Member `gorm:"foreignKey:SpaceID;References:ID;unique;not null"`
	// TODO add reference to marketplace

}

type ShieldAppDomainMapping struct {
	//ID           uint
	ID         string `gorm:"primary_key; size:255"`
	OwnerAppID string
	Url        string
	OwnerApp   ShieldApp `gorm:"foreignKey:OwnerAppID;References:AppId"`
}

// Model for Permissions Table
type Permission struct {
	PermissionId   string `gorm:"primary_key; size:255"`
	PermissionName string `gorm:"size:100; not null; unique"`
	Description    string `gorm:"size:255"`
	Category       string `gorm:"size:255"`
	Mandatory      bool   //to identify default mandatory permissions
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Model for AppPermissions Table
type AppPermission struct {
	AppId        string `gorm:"primary_key; size:255; not null"`
	PermissionId string `gorm:"primary_key; size:255; not null"`
	Mandatory    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Model for AppUserPermissions Table
type AppUserPermission struct {
	UserId       string `gorm:"primary_key; size:255; not null"`
	AppId        string `gorm:"primary_key; size:255; not null"`
	PermissionId string `gorm:"primary_key; size:255; not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Model for UserProvider Table
type UserProvider struct {
	UserId    string `gorm:"primary_key; size:255; not null"`
	Provider  int    `gorm:"primary_key; not null; default:1"` // 1 - shield, 2 - google
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AcPolGrpSub struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	OwnerSpaceID string // FK to Spaces table
	RoleID       string `gorm:"default:null"` // FK to Role table
	OwnerTeamID  string `gorm:"default:null"` // FK to Team table
	OwnerUserID  string `gorm:"default:null"` // FK to User table
	AcPolGrpID   string // FK to AcPolGrp table
	OptCounter   int    `gorm:"size:8"`
	PermissionID string `gorm:"default:null"` // FK to Role table

	AcPermission AcPermissions `gorm:"foreignKey:PermissionID;References:ID"`
	AcPolGrp     AcPolGrp      `gorm:"foreignKey:AcPolGrpID;References:ID"`
	OwnerSpace   Space         `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
	OwnerRole    Role          `gorm:"foreignKey:RoleID;References:ID"`
	OwnerTeam    Team          `gorm:"foreignKey:OwnerTeamID;References:TeamID"`
	OwnerUser    User          `gorm:"foreignKey:OwnerUserID;References:UserID"`
}

type Role struct {
	gorm.Model

	ID           string `gorm:"primaryKey;not null"`
	Name         string `gorm:"index:rolename_unique_index,unique"`
	Description  string
	OwnerSpaceID string `gorm:"index:rolename_unique_index,unique"`
	IsOwner      bool
	CreatedBy    string
	UpdatedBy    string
	OptCounter   int `gorm:"size:8"`

	OwnerSpace Space `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
}

type Team struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TeamID      string `gorm:"primaryKey;not null"`                          // PK and FK to MEMBER table (type T).
	OwnerID     string `gorm:"not null; index:teamname_unique_index,unique"` // PK and FK to MEMBER table  (type S).
	Name        string `gorm:"index:teamname_unique_index,unique"`
	Description string
	Update      string
	UpdatedBy   string
	OptCounter  int `gorm:"size:8"`

	Member     Member `gorm:"foreignKey:TeamID;References:ID"`
	OwnerSpace Space  `gorm:"foreignKey:OwnerID;References:SpaceID"`
}

type AcPermissions struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	Description  string `gorm:"size:255; not null;"`
	Name         string
	OptCounter   int `gorm:"size:8"`
	IsPredefined bool
	Type         int // 1 for internal apps and 2 for consumer apps
	DisplayName  string
}

type PerPolGrps struct {
	gorm.Model
	ID             string `gorm:"primaryKey;not null"`
	AcPermissionID string
	AcPolGrpID     string

	AcPolGrp     AcPolGrp      `gorm:"foreignKey:AcPolGrpID;References:ID"`
	AcPermission AcPermissions `gorm:"foreignKey:AcPermissionID;References:ID"`
}

type MemberRole struct {
	gorm.Model

	ID           string `gorm:"primaryKey;not null"`
	OwnerUserID  string `gorm:"index:member_role_unique_index,unique"` // PK and FK to MEMBER table (type U).
	RoleID       string `gorm:"index:member_role_unique_index,unique"` // NOT NULL	FK to Role.
	OwnerSpaceID string `gorm:"index:member_role_unique_index,unique"` // PK and FK to MEMBER table (type U).
	OptCounter   int    `gorm:"size:8"`

	MemberUser  User  `gorm:"foreignKey:OwnerUserID;references:UserID"`
	Role        Role  `gorm:"foreignKey:RoleID;references:ID"`
	SpaceMember Space `gorm:"foreignKey:OwnerSpaceID;references:SpaceID"`
}

type SpaceMember struct {
	gorm.Model

	ID           string `gorm:"primaryKey;not null"`
	OwnerUserID  string `gorm:"index:member_role_unique_index,unique"` // PK and FK to MEMBER table (type U).
	OwnerSpaceID string `gorm:"index:member_role_unique_index,unique"` // PK and FK to MEMBER table (type U).
	OptCounter   int    `gorm:"size:8"`

	MemberUser  User  `gorm:"foreignKey:OwnerUserID;references:UserID"`
	SpaceMember Space `gorm:"foreignKey:OwnerSpaceID;references:SpaceID"`
}

type TeamMember struct {
	gorm.Model
	ID          string `gorm:"primaryKey;not null"`
	OwnerTeamID string `gorm:"index:team_member_unique_index,unique"` // PK and FK to Team.
	MemberID    string `gorm:"index:team_member_unique_index,unique"` // PK and FK to MEMBER table (type U).
	IsOwner     bool   `gorm:"default:false"`
	OptCounter  int    `gorm:"size:8"`

	Team       Team `gorm:"foreignKey:OwnerTeamID;References:TeamID"`
	UserMember User `gorm:"foreignKey:MemberID;References:UserID"`
}

type DefaultUserSpace struct {
	gorm.Model

	ID           string `gorm:"primaryKey"`
	OwnerUserID  string `gorm:"unique"` // FK to User table
	OwnerSpaceID string // FK to Spaces table
	OptCounter   int    `gorm:"size:8"`
	OwnerUser    User   `gorm:"foreignKey:OwnerUserID;References:UserID"`
	OwnerSpace   Space  `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
}

type AcResource struct {
	gorm.Model
	ID              string `gorm:"primaryKey;not null"`
	Name            string
	Description     string
	Path            string
	FunctionName    string
	EntityName      string
	FunctionMethod  string
	Version         string
	OptCounter      int `gorm:"size:8"`
	OwnerAppID      string
	IsAuthorised    int
	IsAuthenticated int
	OwnerApp        ShieldApp `gorm:"foreignKey:OwnerAppID;References:AppId"`
}

type AcResGrp struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	OwnerSpaceID string // FK to spaces table
	Name         string
	Description  string
	IsPredefined bool
	OptCounter   int `gorm:"size:8"`
	Type         int // 1 for internal apps and 2 for consumer apps

	OwnerSpace Space `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
}

type AcResGpRes struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	AcResGrpID   string
	AcResourceID string // FK to MEMBER table
	OptCounter   int    `gorm:"size:8"`

	AcResource AcResource `gorm:"foreignKey:AcResourceID;references:ID"`
	AcResGrp   AcResGrp   `gorm:"foreignKey:AcResGrpID;references:ID"`
}

type AcResAction struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	AcActionID   string
	AcResourceID string // FK to MEMBER table
	OptCounter   int    `gorm:"size:8"`

	AcResource AcResource `gorm:"foreignKey:AcResourceID;References:ID"`
	AcAction   AcAction   `gorm:"foreignKey:AcActionID;References:ID"`
}

type AcAction struct {
	gorm.Model
	ID          string `gorm:"primaryKey;not null"`
	Name        string
	Description string
	OptCounter  int `gorm:"size:8"`
	OwnerAppID  string
	OwnerApp    ShieldApp `gorm:"foreignKey:OwnerAppID;References:AppId"`
}

type AcActGrp struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	OwnerSpaceID string // FK to spaces table
	Description  string
	Name         string
	IsPredefined bool
	OptCounter   int `gorm:"size:8"`
	Type         int // 1 for internal apps and 2 for consumer apps

	OwnerSpace Space `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
}

type ActGpAction struct {
	gorm.Model
	ID         string `gorm:"primaryKey"`
	AcActGrpID string
	AcActionID string
	OptCounter int `gorm:"size:8"`

	AcAction AcAction `gorm:"foreignKey:AcActionID;References:ID"`
	AcActGrp AcActGrp `gorm:"foreignKey:AcActGrpID;References:ID"`
}

type AcPolicy struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	AcActGrpID   string
	AcResGrpID   string
	OwnerSpaceID string // FK to MEMBER table
	CreatedBy    string // FK to USERS TABLE
	UpdatedBy    string // FK TO USERS TABLE
	Name         string
	Description  string
	Path         string
	OptCounter   int `gorm:"size:8"`
	IsPredefined bool
	Type         int // 1 for internal apps and 2 for consumer apps

	AcActionGroup   AcActGrp `gorm:"foreignKey:AcActGrpID;References:ID"`
	AcResourceGroup AcResGrp `gorm:"foreignKey:AcResGrpID;References:ID"`
	OwnerSpace      Space    `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
	CreatedUser     User     `gorm:"foreignKey:CreatedBy;References:UserID"`
	UpdatedUser     User     `gorm:"foreignKey:UpdatedBy;References:UserID"`
}
type AcPolGrp struct {
	gorm.Model
	ID           string `gorm:"primaryKey;not null"`
	OwnerSpaceID string // FK to MEMBER table
	Description  string `gorm:"size:255; not null;"`
	Name         string
	OptCounter   int `gorm:"size:8"`
	IsPredefined bool
	Type         int // 1 for internal apps and 2 for consumer apps
	EntityType   int // 0 for non entity based policies 1 for block, 2 for app, 3 for environment
	DisplayName  string
	EntityTypes  pq.Int64Array `gorm:"type:integer[]"`
	// type changed to array and renamed to types  0 for non entity based policies 1 for block, 2 for app, 3 for environment

	OwnerSpace Space `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
}

// Model for EntityDefinition Table
type EntityTypeDefinition struct {
	ID          int    `gorm:"primaryKey;not null;autoIncrement"`
	Name        string `gorm:"unique"`
	DisplayName string `gorm:"unique"`
}

type Entities struct {
	EntityID   string `gorm:"primaryKey;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Type       int64
	Definition EntityTypeDefinition `gorm:"foreignKey:Type;references:ID"`
	Label      string
}

type EntitySpaceMapping struct {
	gorm.Model

	ID            string `gorm:"primaryKey;not null"`
	OwnerEntityID string
	OwnerSpaceID  string

	OwnerSpace Space    `gorm:"foreignKey:OwnerSpaceID;References:SpaceID"`
	Entities   Entities `gorm:"foreignKey:OwnerEntityID;References:EntityID"`
}

type PolGrpSubsEntityMapping struct {
	gorm.Model

	ID            string `gorm:"primaryKey;not null"`
	OwnerEntityID string
	PolGrpSubsID  string

	PolicyGroupSubs AcPolGrpSub `gorm:"foreignKey:PolGrpSubsID;References:ID"`
	Entities        Entities    `gorm:"foreignKey:OwnerEntityID;References:EntityID"`
}

type PredefinedEnv struct {
	gorm.Model
	ID   string `gorm:"primaryKey; not null"`
	Name string
}

type PolGpPolicy struct {
	gorm.Model
	ID         string `gorm:"primaryKey;not null"`
	AcPolicyID string
	AcPolGrpID string // FK to MEMBER table
	OptCounter int    `gorm:"size:8"`

	AcPolGrp AcPolGrp `gorm:"foreignKey:AcPolGrpID;References:ID"`
	AcPolicy AcPolicy `gorm:"foreignKey:AcPolicyID;References:ID"`
}

type Invites struct {
	gorm.Model
	ID         string `gorm:"primaryKey;not null"`
	Notes      string
	CreatedBy  string
	Status     int // 1-pending 2 complete 3 declined
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	InviteType int //1-email invyt 2 invite link
	InviteLink string
	InviteCode string
	Email      string

	CreatedUser User `gorm:"foreignKey:CreatedBy;References:UserID"`
}

type InviteDetails struct {
	gorm.Model
	ID             string `gorm:"primaryKey;not null"`
	InviteID       string
	InvitedSpaceID string
	InvitedTeamID  string `gorm:"default: null"`
	InvitedRoleID  string `gorm:"default: null"`
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Invite Invites `gorm:"foreignKey:InviteID;References:ID"`
	Space  Space   `gorm:"foreignKey:InvitedSpaceID;References:SpaceID"`
	Team   Team    `gorm:"foreignKey:InvitedTeamID;References:TeamID"`
	Role   Role    `gorm:"foreignKey:InvitedRoleID;References:ID"`
}
