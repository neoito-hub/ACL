package roles_list_permissions

import (
	"time"

	"gorm.io/datatypes"
)

type RequestObject struct {
	RoleID        string `json:"role_id"`
	SearchKeyword string `json:"search_keyword"`
	PageLimit     int    `json:"page_limit"`
	Offset        int    `json:"offset"`
	SortColumn    string `json:"sort_column"`
	SortDirection string `json:"sort_direction"`
}

type Response struct {
	Err  bool        `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResultData struct {
	TotalCount int                `json:"total_count"`
	Data       []PoliciesListData `json:"data"`
}

type PoliciesListData struct {
	PermissionID     string           `json:"permission_id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	PolicyGroups     datatypes.JSON   `json:"policy_groups"`
	PgCount          int              `json:"pg_count"`
	Entities         datatypes.JSON   `json:"entities"`
	IsPredefined     bool             `json:"is_predefined"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	AttachedEntities AttachedEntities `json:"attached_entities"`
	EntityTypes      datatypes.JSON   `json:"entity_types"`
}

type AttachedEntities struct {
	SpaceAccessEntities Entities
	AddedEntities       Entities
}

type Entities map[int][]Entity

type Entity struct {
	EntityID   string `json:"entity_id"`
	EntityType int    `json:"entity_type"`
	Label      string `json:"label"`
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
