package teams_list_available_permissions

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	db := payload.Db

	var policiesList []PoliciesListData

	var filterString string

	valuesMap := make(map[string]interface{})
	availablePermissionQuery := `
	with acl_pgs as (
		select pg.id,pg.display_name,pg.name,pg.entity_types from
		ac_pol_grps pg inner join
		ac_pol_grp_subs pgs ON pg.id = pgs.ac_pol_grp_id and  pgs.owner_space_id=@owner_space_id
		LEFT JOIN (select mr.role_id, r.name,r.is_owner  from member_roles mr
			inner join roles r on r.id = mr.role_id
			where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id
		   ) r on r.role_id = pgs.role_id
 LEFT JOIN (select t.team_id, name from team_members tm 
			inner join teams t on tm.owner_team_id=t.team_id 
			where tm.member_id=@owner_user_id and t.owner_id=@owner_space_id
		   ) t ON t.team_id = pgs.owner_team_id
		   WHERE  
		  pgs.owner_user_id=@owner_user_id OR r.role_id is not null or t.team_id is not null
	),
	owner_pgs as (
		select pg.id,pg.display_name,pg.name,pg.entity_types from
		ac_pol_grps pg 
		LEFT JOIN (select mr.role_id, r.name,r.is_owner  from member_roles mr
			inner join roles r on r.id = mr.role_id
			where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id and r.is_owner
		   ) r on true
		   WHERE  
		  r.is_owner and pg.is_predefined and pg.type=1
	),
	selected_user_pgs as (
		select distinct pgs.permission_id from
		ac_pol_grp_subs pgs 
		INNER JOIN ac_pol_grps pg ON pg.id = pgs.ac_pol_grp_id
		INNER JOIN teams t ON t.team_id = pgs.owner_team_id 
		WHERE  pgs.owner_space_id=@owner_space_id AND t.team_id = @selected_team_id
	)

	SELECT acper.id as permission_id,acper.display_name as name, acper.description,json_agg(distinct jsonb_build_object('policy_group_id',pg.id,'policy_group_display_name',pg.display_name,'policy_group_name',pg.name)) policy_groups,count(pg.id) as pg_count,json_object_agg(distinct coalesce(pg.entity_type,0),true) as entity_types FROM 
	(select acper.* from ac_permissions acper left join selected_user_pgs selpgs on selpgs.permission_id=acper.id where selpgs.permission_id is null )acper
	left join per_pol_grps prpgp on prpgp.ac_permission_id=acper.id
	left join (
		select acl_pgs.*,unnest(acl_pgs.entity_types) as entity_type from acl_pgs union 
		select owner_pgs.*,unnest(owner_pgs.entity_types)as entity_type from
		owner_pgs
	)pg on pg.id=prpgp.ac_pol_grp_id
	
	 `

	filterString = ``

	availabePermissionCountQuery := `
	with acl_pgs as (
		select pg.id,pg.display_name,pg.name,pg.entity_type from
		ac_pol_grps pg inner join
		ac_pol_grp_subs pgs ON pg.id = pgs.ac_pol_grp_id and  pgs.owner_space_id=@owner_space_id
		LEFT JOIN (select mr.role_id, r.name,r.is_owner  from member_roles mr
			inner join roles r on r.id = mr.role_id
			where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id
		   ) r on r.role_id = pgs.role_id
 LEFT JOIN (select t.team_id, name from team_members tm 
			inner join teams t on tm.owner_team_id=t.team_id 
			where tm.member_id=@owner_user_id and t.owner_id=@owner_space_id
		   ) t ON t.team_id = pgs.owner_team_id
		   WHERE  
		  pgs.owner_user_id=@owner_user_id OR r.role_id is not null or t.team_id is not null
	),
	owner_pgs as (
		select pg.id,pg.display_name,pg.name,pg.entity_type from
		ac_pol_grps pg 
		LEFT JOIN (select mr.role_id, r.name,r.is_owner  from member_roles mr
			inner join roles r on r.id = mr.role_id
			where mr.owner_user_id=@owner_user_id and mr.owner_space_id=@owner_space_id and r.is_owner
		   ) r on true
		   WHERE  
		  r.is_owner and pg.is_predefined and pg.type=1
	),
	selected_user_pgs as (
		select distinct pgs.permission_id from
		ac_pol_grp_subs pgs 
		INNER JOIN ac_pol_grps pg ON pg.id = pgs.ac_pol_grp_id
		INNER JOIN teams t ON t.team_id = pgs.owner_team_id 
		WHERE  pgs.owner_space_id=@owner_space_id AND t.team_id = @selected_team_id
	)

	SELECT count(distinct acper.id) from
	(select acper.* from ac_permissions acper left join selected_user_pgs selpgs on selpgs.permission_id=acper.id where selpgs.permission_id is null )acper
	left join per_pol_grps prpgp on prpgp.ac_permission_id=acper.id
	left join (select acl_pgs.* from acl_pgs union
	select owner_pgs.* from
	owner_pgs
	)pg on pg.id=prpgp.ac_pol_grp_id `

	if len(b.SearchKeyword) > 0 {
		//for the attached permissions
		searchFilter := ""
		searchFilter = "acper.display_name ilike @Keyword"

		valuesMap["Keyword"] = "%" + b.SearchKeyword + "%"

		AttachToMainfilter(&searchFilter, &filterString)

	}

	SortColumns := make(map[string]string)
	SortDirections := make(map[string]string)
	if len(b.SortColumn) == 0 {
		b.SortColumn = "PolicyCount"
	}
	if len(b.SortDirection) == 0 {
		b.SortDirection = "desc"
	}

	SortColumns[b.SortColumn] = b.SortColumn
	SortDirections[b.SortDirection] = b.SortDirection

	SortColumns["PolicyCount"] = "pg_count"
	SortColumns["updatedAt"] = "updated_at"
	SortDirections["desc"] = "desc"
	SortDirections["asc"] = "asc"

	orderByString := ` order by ` + SortColumns[b.SortColumn] + " " + SortDirections[strings.ToLower(b.SortDirection)]

	availablePermissionQuery = availablePermissionQuery + filterString + " group by acper.id,acper.display_name,acper.description " + orderByString + " limit @limit offset @offset"
	availabePermissionCountQuery = availabePermissionCountQuery + filterString
	valuesMap["limit"] = b.PageLimit
	valuesMap["offset"] = b.Offset
	valuesMap["owner_user_id"] = payload.UserID
	valuesMap["owner_space_id"] = payload.SpaceID
	valuesMap["selected_team_id"] = b.TeamID

	res := db.Raw(availablePermissionQuery, valuesMap).Scan(&policiesList)

	if res.Error != nil {

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	var resultData ResultData

	countRes := db.Raw(availabePermissionCountQuery, valuesMap).Scan(&resultData.TotalCount)

	if countRes.Error != nil {

		handlerResp = common_services.BuildErrorResponse(true, "Invalid request payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	resultData.Data = policiesList

	handlerResp = common_services.BuildResponse(false, "Available permissions for team listed successfully!", Response{Data: resultData, Err: false, Msg: "Available permissions for team listed successfully!"}, http.StatusOK)
	return handlerResp
}
