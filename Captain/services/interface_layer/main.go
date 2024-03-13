package interface_layer

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func Handler(shieldUser string, req *http.Request, writer http.ResponseWriter, db *gorm.DB, actionName string, functionName string, spaceID string) (OwnerCheckData, error) {

	// var appDetails AppData
	var ownerCheckData OwnerCheckData

	ownerCheckQuery := `select exists(select mr.role_id from member_roles mr 
		inner join roles r on mr.role_id=r.id 
		where mr.owner_space_id=? and mr.owner_user_id=? and r.is_owner)`

	db.Raw(ownerCheckQuery, spaceID, shieldUser).Scan(&ownerCheckData)

	if !ownerCheckData.Exists {

		var aclDetails AclData

		aclQuery := `
	with
	
	--getting roles for the user via member_roles table for the given space
	user_roles as (select mr.role_id from member_roles mr where mr.owner_space_id=@SpaceID 
				   and mr.owner_user_id=@UserID),
	
	--getting teams for the user via team_members table for the given space
	user_teams as (select t.team_id from team_members tm inner join teams t on tm.owner_team_id=t.team_id where t.owner_id=@SpaceID
		 and tm.member_id=@UserID)
		
	select exists(select acr.id,acr.name,act.id,act.name from 
	--policy group filtering for the given user
	(select ac_pol_grp_id as id from ac_pol_grp_subs grp 
	 
	 --joining with user roles
	 left join user_roles ur on ur.role_id=grp.role_id 
	 
	 --joining with user teams
	 left join user_teams ut on ut.team_id=grp.owner_team_id 
	 where owner_space_id=@SpaceID  and (grp.owner_user_id=@UserID 
	or ur.role_id is not null or ut.team_id is not null)
	group by grp.ac_pol_grp_id
	) grp
	
	--joining with policies via pol_gp_policies bridge table
	inner join pol_gp_policies pgbr on pgbr.ac_pol_grp_id=grp.id 
	inner join ac_policies pol on pol.id=pgbr.ac_policy_id
	
	
	--joining with resource groups and action groups
	inner join ac_res_grps resgrp on resgrp.id=pol.ac_res_grp_id
	inner join ac_act_grps acgrp on acgrp.id=pol.ac_act_grp_id
	
	--joining with bridge table
	inner join ac_res_gp_res resgrpbr on resgrpbr.ac_res_grp_id=resgrp.id
	inner join act_gp_actions actgrpbr on actgrpbr.ac_act_grp_id=acgrp.id
	
	-- joining with actions and resources
	inner join ac_resources acr on (acr.id=resgrpbr.ac_resource_id and acr.function_name=@FunctionName) 
	inner join ac_actions act on (act.id=actgrpbr.ac_action_id and act.name=@ActionName ))`

		valuesMap := make(map[string]interface{})

		valuesMap["ActionName"] = actionName
		valuesMap["FunctionName"] = functionName
		valuesMap["UserID"] = shieldUser
		valuesMap["SpaceID"] = spaceID

		// TODO
		aclRes := db.Raw(aclQuery, valuesMap).Scan(&aclDetails)

		if aclRes.Error != nil {

			return ownerCheckData, errors.New("resource access forbidden")
		}

		if aclRes.RowsAffected < 1 || !aclDetails.Exists {

			return ownerCheckData, errors.New("resource access forbidden")
		}
		fmt.Println(aclDetails.Exists)

		if !aclDetails.Exists {
			return ownerCheckData, errors.New("resource access forbidden")
		}
	}

	return ownerCheckData, nil

}
