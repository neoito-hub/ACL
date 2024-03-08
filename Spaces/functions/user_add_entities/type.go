package user_add_entities

import "gorm.io/datatypes"

type RequestObject struct {
	NewEntityMappings     []NewEntityMappings `json:"new_entity_mappings"`
	DeletedEntityMappings []string            `json:"deleted_entity_mappings"`
	UserID                string              `json:"user_id"`
	TypeID                int                 `json:"type_id"`
}

type NewEntityMappings struct {
	EntityID   string `json:"entity_id"`
	AcPolGrpID string `json:"ac_pol_grp_id"`
}

type ExistingMappings struct {
	Polgrpsubs datatypes.JSON `json:"polgrpsubs"`
	Etmappings datatypes.JSON `json:"etmappings"`
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
