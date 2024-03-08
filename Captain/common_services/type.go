package common_services

import "gorm.io/gorm"

type HandlerPayload struct {
	Url         string
	RequestBody string
	UserID      string
	Db          *gorm.DB
	UserName    string
	Queryparams map[string]string
}

type HandlerResponse struct {
	Err    bool   `json:"err"`
	Data   string `json:"data"`
	Status int    `json:"status"`
}

type DBInfo struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	Sslmode  string
	Timezone string
}

type Resources struct {
	EntiyName       string `json:"entity_name"`
	IsAuthorised    int    `json:"is_authorised"`
	IsAuthenticated int    `json:"is_authenticated"`
}

type ContextData struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	SpaceID  string `json:"space_id"`
	IsOwner  bool   `json:"is_owner"`
}
