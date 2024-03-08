package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/neoito-hub/ACL-Block/spaces/functions/accept_invite"
	"github.com/neoito-hub/ACL-Block/spaces/functions/cancel_invite"
	"github.com/neoito-hub/ACL-Block/spaces/functions/check_assigned_block_to_app"
	"github.com/neoito-hub/ACL-Block/spaces/functions/check_business_name"
	"github.com/neoito-hub/ACL-Block/spaces/functions/check_entity_name"
	"github.com/neoito-hub/ACL-Block/spaces/functions/check_role_name"
	"github.com/neoito-hub/ACL-Block/spaces/functions/check_space_name"
	"github.com/neoito-hub/ACL-Block/spaces/functions/check_team_name"
	"github.com/neoito-hub/ACL-Block/spaces/functions/create_entity"
	"github.com/neoito-hub/ACL-Block/spaces/functions/create_invite_link"
	"github.com/neoito-hub/ACL-Block/spaces/functions/create_logo_signed_url"
	"github.com/neoito-hub/ACL-Block/spaces/functions/create_role"
	"github.com/neoito-hub/ACL-Block/spaces/functions/create_space"
	"github.com/neoito-hub/ACL-Block/spaces/functions/create_team"
	"github.com/neoito-hub/ACL-Block/spaces/functions/delete_space"
	"github.com/neoito-hub/ACL-Block/spaces/functions/delete_team"
	"github.com/neoito-hub/ACL-Block/spaces/functions/get_invite_by_id"
	"github.com/neoito-hub/ACL-Block/spaces/functions/get_space_by_id"
	"github.com/neoito-hub/ACL-Block/spaces/functions/get_user_by_id"
	"github.com/neoito-hub/ACL-Block/spaces/functions/invite_user_email"
	"github.com/neoito-hub/ACL-Block/spaces/functions/list_invited_users"
	"github.com/neoito-hub/ACL-Block/spaces/functions/list_roles"
	"github.com/neoito-hub/ACL-Block/spaces/functions/list_spaces"
	"github.com/neoito-hub/ACL-Block/spaces/functions/list_spaces_detailed"
	"github.com/neoito-hub/ACL-Block/spaces/functions/list_teams"
	"github.com/neoito-hub/ACL-Block/spaces/functions/list_users"
	"github.com/neoito-hub/ACL-Block/spaces/functions/resend_invite_email"
	"github.com/neoito-hub/ACL-Block/spaces/functions/revoke_invite"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_add_entities"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_add_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_add_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_create_invite_link"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_delete_existing_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_delete_user"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_invite_user_email"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_list_available_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_list_entities"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_list_existing_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_list_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_list_to_add_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_list_users"
	"github.com/neoito-hub/ACL-Block/spaces/functions/roles_search_user"
	"github.com/neoito-hub/ACL-Block/spaces/functions/search_user"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_add_entities"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_add_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_add_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_create_invite_link"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_delete_existing_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_delete_user"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_invite_user_email"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_list_available_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_list_entities"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_list_existing_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_list_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_list_to_add_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_list_users"
	"github.com/neoito-hub/ACL-Block/spaces/functions/teams_search_user"
	"github.com/neoito-hub/ACL-Block/spaces/functions/unassign_block_from_app"
	"github.com/neoito-hub/ACL-Block/spaces/functions/update_role"
	"github.com/neoito-hub/ACL-Block/spaces/functions/update_space"
	"github.com/neoito-hub/ACL-Block/spaces/functions/update_team"
	"github.com/neoito-hub/ACL-Block/spaces/functions/update_user"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_add_entities"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_add_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_add_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_delete_existing_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_available_entities"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_available_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_entities"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_existing_pol_grp_subs"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_permissions"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_pol_grp_subs_from_roles"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_pol_grp_subs_from_teams"
	"github.com/neoito-hub/ACL-Block/spaces/functions/user_list_to_add_pol_grp_subs"

	"github.com/neoito-hub/ACL-Block/spaces/functions/list_entity_definition"

	"github.com/aidarkhanov/nanoid"
	"github.com/neoito-hub/ACL-Block/spaces/common_services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var funcs map[string]interface{}
var db *gorm.DB

func InvokeSpacesFunction(funcs map[string]interface{}, payload common_services.HandlerPayload) common_services.HandlerResponse {

	f, functionExists := funcs[payload.Url]
	if !functionExists {
		return common_services.HandlerResponse{Data: "", Status: 404, Err: true}
	}

	timed := InvokeGrpcFunction(f).(func(common_services.HandlerPayload) common_services.HandlerResponse)

	result := timed(common_services.HandlerPayload{Url: payload.Url, UserID: payload.UserID, RequestBody: payload.RequestBody, Db: db, Queryparams: payload.Queryparams, UserName: payload.UserName, SpaceID: payload.SpaceID, IsOwner: payload.IsOwner})

	return result
}

func loadFuncs() error {
	funcs = map[string]interface{}{
		"/api/spaces/v0.1/create-logo-signed-url/invoke":             create_logo_signed_url.Handler,
		"/api/spaces/v0.1/create-space/invoke":                       create_space.Handler,
		"/api/spaces/v0.1/update-space/invoke":                       update_space.Handler,
		"/api/spaces/v0.1/search-user/invoke":                        search_user.Handler,
		"/api/spaces/v0.1/teams-search-user/invoke":                  teams_search_user.Handler,
		"/api/spaces/v0.1/roles-search-user/invoke":                  roles_search_user.Handler,
		"/api/spaces/v0.1/get-space-by-id/invoke":                    get_space_by_id.Handler,
		"/api/spaces/v0.1/accept-invite/invoke":                      accept_invite.Handler,
		"/api/spaces/v0.1/delete-space/invoke":                       delete_space.Handler,
		"/api/spaces/v0.1/list-spaces/invoke":                        list_spaces.Handler,
		"/api/spaces/v0.1/list-spaces-detailed/invoke":               list_spaces_detailed.Handler,
		"/api/spaces/v0.1/check-space-name/invoke":                   check_space_name.Handler,
		"/api/spaces/v0.1/check-entity-name/invoke":                  check_entity_name.Handler,
		"/api/spaces/v0.1/check-business-name/invoke":                check_business_name.Handler,
		"/api/spaces/v0.1/create-invite-link/invoke":                 create_invite_link.Handler,
		"/api/spaces/v0.1/send-user-invite-email/invoke":             invite_user_email.Handler,
		"/api/spaces/v0.1/teams-create-invite-link/invoke":           teams_create_invite_link.Handler,
		"/api/spaces/v0.1/teams-send-user-invite-email/invoke":       teams_invite_user_email.Handler,
		"/api/spaces/v0.1/roles-send-user-invite-email/invoke":       roles_invite_user_email.Handler,
		"/api/spaces/v0.1/roles-create-invite-link/invoke":           roles_create_invite_link.Handler,
		"/api/spaces/v0.1/list-users/invoke":                         list_users.Handler,
		"/api/spaces/v0.1/list-invited-users/invoke":                 list_invited_users.Handler,
		"/api/spaces/v0.1/teams-list-users/invoke":                   teams_list_users.Handler,
		"/api/spaces/v0.1/roles-list-users/invoke":                   roles_list_users.Handler,
		"/api/spaces/v0.1/create-team/invoke":                        create_team.Handler,
		"/api/spaces/v0.1/update-team/invoke":                        update_team.Handler,
		"/api/spaces/v0.1/delete-team/invoke":                        delete_team.Handler,
		"/api/spaces/v0.1/list-teams/invoke":                         list_teams.Handler,
		"/api/spaces/v0.1/check-team-name/invoke":                    check_team_name.Handler,
		"/api/spaces/v0.1/create-role/invoke":                        create_role.Handler,
		"/api/spaces/v0.1/list-roles/invoke":                         list_roles.Handler,
		"/api/spaces/v0.1/update-role/invoke":                        update_role.Handler,
		"/api/spaces/v0.1/check-role-name/invoke":                    check_role_name.Handler,
		"/api/spaces/v0.1/revoke-invite/invoke":                      revoke_invite.Handler,
		"/api/spaces/v0.1/resend-invite-email/invoke":                resend_invite_email.Handler,
		"/api/spaces/v0.1/get-invite-by-id/invoke":                   get_invite_by_id.Handler,
		"/api/spaces/v0.1/cancel-invite/invoke":                      cancel_invite.Handler,
		"/api/spaces/v0.1/get-user-by-id/invoke":                     get_user_by_id.Handler,
		"/api/spaces/v0.1/update-user/invoke":                        update_user.Handler,
		"/api/spaces/v0.1/teams-delete-user/invoke":                  teams_delete_user.Handler,
		"/api/spaces/v0.1/roles-delete-user/invoke":                  roles_delete_user.Handler,
		"/api/spaces/v0.1/roles-list-existing-pol-grp-subs/invoke":   roles_list_existing_pol_grp_subs.Handler,
		"/api/spaces/v0.1/roles-delete-existing-pol-grp-subs/invoke": roles_delete_existing_pol_grp_subs.Handler,
		"/api/spaces/v0.1/roles-list-to-add-pol-grp-subs/invoke":     roles_list_to_add_pol_grp_subs.Handler,
		"/api/spaces/v0.1/roles-add-pol-grp-subs/invoke":             roles_add_pol_grp_subs.Handler,
		"/api/spaces/v0.1/teams-list-existing-pol-grp-subs/invoke":   teams_list_existing_pol_grp_subs.Handler,
		"/api/spaces/v0.1/teams-delete-existing-pol-grp-subs/invoke": teams_delete_existing_pol_grp_subs.Handler,
		"/api/spaces/v0.1/teams-list-to-add-pol-grp-subs/invoke":     teams_list_to_add_pol_grp_subs.Handler,
		"/api/spaces/v0.1/teams-add-pol-grp-subs/invoke":             teams_add_pol_grp_subs.Handler,
		"/api/spaces/v0.1/user-list-existing-pol-grp-subs/invoke":    user_list_existing_pol_grp_subs.Handler,
		"/api/spaces/v0.1/user-list-pol-grp-subs-from-teams/invoke":  user_list_pol_grp_subs_from_teams.Handler,
		"/api/spaces/v0.1/user-list-pol-grp-subs-from-roles/invoke":  user_list_pol_grp_subs_from_roles.Handler,
		"/api/spaces/v0.1/user-delete-existing-pol-grp-subs/invoke":  user_delete_existing_pol_grp_subs.Handler,
		"/api/spaces/v0.1/user-list-to-add-pol-grp-subs/invoke":      user_list_to_add_pol_grp_subs.Handler,
		"/api/spaces/v0.1/user-add-pol-grp-subs/invoke":              user_add_pol_grp_subs.Handler,
		"/api/spaces/v0.1/user-list-entities/invoke":                 user_list_entities.Handler,
		"/api/spaces/v0.1/user-add-entities/invoke":                  user_add_entities.Handler,
		"/api/spaces/v0.1/roles-list-entities/invoke":                roles_list_entities.Handler,
		"/api/spaces/v0.1/roles-add-entities/invoke":                 roles_add_entities.Handler,
		"/api/spaces/v0.1/teams-list-entities/invoke":                teams_list_entities.Handler,
		"/api/spaces/v0.1/teams-add-entities/invoke":                 teams_add_entities.Handler,

		"/api/spaces/v0.1/unassign-block-from-app/invoke":          unassign_block_from_app.Handler,
		"/api/spaces/v0.1/check-assigned-block-to-app/invoke":      check_assigned_block_to_app.Handler,
		"/api/spaces/v0.1/user-list-permissions/invoke":            user_list_permissions.Handler,
		"/api/spaces/v0.1/user-list-available-permissions/invoke":  user_list_available_permissions.Handler,
		"/api/spaces/v0.1/user-list-available-entities/invoke":     user_list_available_entities.Handler,
		"/api/spaces/v0.1/user-add-permissions/invoke":             user_add_permissions.Handler,
		"/api/spaces/v0.1/teams-list-permissions/invoke":           teams_list_permissions.Handler,
		"/api/spaces/v0.1/teams-list-available-permissions/invoke": teams_list_available_permissions.Handler,
		"/api/spaces/v0.1/teams-add-permissions/invoke":            teams_add_permissions.Handler,
		"/api/spaces/v0.1/roles-list-permissions/invoke":           roles_list_permissions.Handler,
		"/api/spaces/v0.1/roles-add-permissions/invoke":            roles_add_permissions.Handler,
		"/api/spaces/v0.1/roles-list-available-permissions/invoke": roles_list_available_permissions.Handler,
		"/api/spaces/v0.1/list-entity-definition/invoke":           list_entity_definition.Handler,
		"/api/spaces/v0.1/create-entity/invoke":                    create_entity.Handler,
	}

	return nil
}

func InvokeGrpcFunction(f interface{}) interface{} {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		panic("expects a function")
	}
	vf := reflect.ValueOf(f)
	wrapperF := reflect.MakeFunc(rf, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := vf.Call(in)
		end := time.Now()
		fmt.Printf("calling %s took %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}

func DBInit() {
	dbinf := &common_services.DBInfo{}
	var dbErr error

	dbinf.Host = os.Getenv("SPACES_POSTGRES_HOST")
	dbinf.User = os.Getenv("SPACES_POSTGRES_USER")
	dbinf.Password = os.Getenv("SPACES_POSTGRES_PASSWORD")
	dbinf.Name = os.Getenv("SPACES_POSTGRES_NAME")
	dbinf.Port = os.Getenv("SPACES_POSTGRES_PORT")
	dbinf.Sslmode = os.Getenv("SPACES_POSTGRES_SSLMODE")
	dbinf.Timezone = os.Getenv("SPACES_POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbinf.Host, dbinf.User, dbinf.Password, dbinf.Name, dbinf.Port, dbinf.Sslmode, dbinf.Timezone)
	db, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if dbErr != nil {
		panic("DB connection err")
	}

	// UpdateAcResources();
}

func UpdateAcResources() {
	entityName := "spaces"
	appName := "JqCFOyekIFZfYo-kP6CA5"
	var newResources []string

	var ownerAppID string

	db.Raw(`select app_id from shield_apps where app_name=?`, appName).Scan(&ownerAppID)

	//closing connection to db
	for path, _ := range funcs {
		functionSlice := strings.Split(path, "/")
		var resourceExists bool
		isAuthorised := 2
		isAuthenticated := 2
		functionMethod := "POST"

		functionName := functionSlice[len(functionSlice)-2]
		version := functionSlice[3]
		fmt.Printf("function name is%v\n", functionName)
		fmt.Printf("version is %v", version)

		db.Raw(`select exists(select * from ac_resources where owner_app_id=? and function_name=?)`, ownerAppID, functionName).Scan(&resourceExists)
		fmt.Printf("resource exists is %v", resourceExists)
		if resourceExists {

			// 	db.Raw(`update ac_resources set function_name=?,path=?,is_authorised=?,is_authenticated=? where function_name=? and
			// owner_app_id=?`, functionName, path, functionName, isAuthorised, isAuthenticated, ownerAppID)

		} else {
			newResources = append(newResources, functionName)
			fmt.Println(isAuthorised, isAuthenticated, functionMethod, functionName)
			db.Exec(`INSERT INTO public.ac_resources(
			id, created_at, updated_at, name, path, function_name, entity_name, function_method, version, owner_app_id, is_authorised, is_authenticated)
			VALUES (?,now(),now(), ?, ?, ?, ?, ?, ?, ?, ?, ?)`, nanoid.New(), fmt.Sprintf("Resource for %v", functionName), path, functionName, entityName, functionMethod, version, ownerAppID, isAuthorised, isAuthenticated)
		}

	}

	newResourceString, _ := json.Marshal(newResources)

	fmt.Printf("new resources to be added are   %v\n", string(newResourceString))

}

func CloseDbCOnn() {
	//closing connection to db
	sqlDB, dberr := db.DB()
	if dberr != nil {
		log.Fatalln(dberr)
	}
	defer sqlDB.Close()

}
