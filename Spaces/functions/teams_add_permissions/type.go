package teams_add_permissions

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type RequestObject struct {
	Permissions []Permission `json:"permissions"`
	TeamID      string       `json:"team_id"`
}

type Permission struct {
	PermissionID               string   `json:"permission_id"`
	AddedEntities              []Entity `json:"added_entities"`
	DeletedEntities            []Entity `json:"deleted_entities"`
	IsDelete                   bool     `json:"is_delete"`
	AddedSpaceAccessEntities   []Entity `json:"added_space_access_entities"`
	DeletedSpaceAccessEntities []Entity `json:"deleted_space_access_entities"`
}

type Entity struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
}

type SpaceAccessEntity struct {
	Type     int  `json:"type"`
	IsDelete bool `json:"is_delete"`
}

type ExistingPermissions struct {
	ID                     string         `json:"id"`
	PolicyGroups           datatypes.JSON `json:"policy_groups"`
	Entities               datatypes.JSON `json:"entities"`
	PredefinedPolicyGroups datatypes.JSON `json:"predefined_policy_groups"`
}

type PolicyGroup struct {
	SubsID      string        `json:"subs_id"`
	PolGrpID    string        `json:"polgrp_id"`
	EntityTypes pq.Int64Array `json:"entity_types"`
	NewSubsID   string        `json:"new_subs_id"`
}

type Entities struct {
	EntityID        string `json:"entity_id"`
	EntityType      int    `json:"entity_type"`
	ChangeType      int    `json:"change_type"` // 1 for add and 2 for delete
	EntityMappingID string `json:"entity_mapping_id"`
}

type PermissionsMetadata struct {
	PolicyGroupsMap           map[string]ExistingPolicyGroups `json:"policy_groups_map"`
	EntitiesMap               map[string]ExistingEntities     `json:"entities_map"`
	PredefinedPolicyGroupsMap map[string]ExistingPolicyGroups `json:"predefined_policy_groups_map"`
}

type ExistingPolicyGroups struct {
	SubsIDs     []string      `json:"subs_ids"`
	PolGrpID    string        `json:"polgrp_id"`
	EntityTypes pq.Int64Array `json:"entity_types"`
	NewSubsID   string        `json:"new_subs_id"`
}

type ExistingEntities struct {
	EntityID         string   `json:"entity_id"`
	EntityType       int      `json:"entity_type"`
	ChangeType       int      `json:"change_type"` // 1 for add and 2 for delete
	EntityMappingIDs []string `json:"entity_mapping_ids"`
}

type ResponseData struct {
	ID         string `json:"id"`
	SpaceID    string `json:"space_id"`
	UserID     string `json:"user_id"`
	AcPolGrpID string `json:"ac_pol_grp_id"`
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
