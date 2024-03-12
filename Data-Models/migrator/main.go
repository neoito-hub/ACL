package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neoito-hub/ACL-Block/Data-Models/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInfo struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	Sslmode  string
	Timezone string
}

func DBInit() *gorm.DB {
	dbinf := &DBInfo{}

	dbinf.Host = os.Getenv("POSTGRES_HOST")
	dbinf.User = os.Getenv("POSTGRES_USER")
	dbinf.Password = os.Getenv("POSTGRES_PASSWORD")
	dbinf.Name = os.Getenv("POSTGRES_NAME")
	dbinf.Port = os.Getenv("POSTGRES_PORT")
	dbinf.Sslmode = os.Getenv("POSTGRES_SSLMODE")
	dbinf.Timezone = os.Getenv("POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbinf.Host, dbinf.User, dbinf.Password, dbinf.Name, dbinf.Port, dbinf.Sslmode, dbinf.Timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("DB connection err")
	}

	return db
}

func main() {
	// Load env vars
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := DBInit()

	// DropTable(db)
	Migrate(db)
}

func Migrate(db *gorm.DB) {

	resNanoId := db.Exec(`
	CREATE OR REPLACE FUNCTION public.nanoid(
		size integer DEFAULT 21)
		RETURNS text
		LANGUAGE 'plpgsql'
		COST 100
		STABLE PARALLEL UNSAFE
	AS $BODY$
	DECLARE
		id text := '';
		i int := 0;
		urlAlphabet char(64) := '_-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
		bytes bytea := gen_random_bytes(size);
		byte int;
		pos int;
	BEGIN
		WHILE i < size LOOP
			byte := get_byte(bytes, i);
			pos := (byte & 63) + 1;
			id := id || substr(urlAlphabet, pos, 1);
			i = i + 1;
		END LOOP;
		RETURN id;
	END
	$BODY$;
	
	ALTER FUNCTION public.nanoid(integer)
		OWNER TO postgres;
	`)

	if resNanoId.Error != nil {
		log.Fatal("Error")
	}

	resRandomBytes := db.Exec(`
	CREATE OR REPLACE FUNCTION public.gen_random_bytes(
		integer)
		RETURNS bytea
		LANGUAGE 'c'
		COST 1
		VOLATILE STRICT PARALLEL SAFE 
	AS '$libdir/pgcrypto', 'pg_random_bytes'
	;
	
	ALTER FUNCTION public.gen_random_bytes(integer)
		OWNER TO postgres;`)

	if resRandomBytes.Error != nil {
		log.Fatal("Error")
	}

	mbrErr := db.AutoMigrate(&models.Member{})
	if mbrErr != nil {
		log.Fatalf("Error AutoMigrate Member %v", mbrErr)
	}

	spErr := db.AutoMigrate(&models.Space{})
	if spErr != nil {

		log.Fatalf("Error AutoMigrate Space %v", spErr)
	}

	usrPrErr := db.AutoMigrate(&models.UserProvider{})
	if usrPrErr != nil {
		log.Fatalf("Error AutoMigrate UserProvider %v", usrPrErr)
	}

	appUsrPerErr := db.AutoMigrate(&models.AppUserPermission{})
	if appUsrPerErr != nil {
		log.Fatalf("Error AutoMigrate AppUserPermission %v", appUsrPerErr)
	}

	appPerErr := db.AutoMigrate(&models.AppPermission{})
	if appPerErr != nil {
		log.Fatalf("Error AutoMigrate AppPermission %v", appPerErr)
	}

	perErr := db.AutoMigrate(&models.Permission{})
	if perErr != nil {
		log.Fatalf("Error AutoMigrate Permission %v", perErr)
	}

	shldErr := db.AutoMigrate(&models.ShieldApp{})
	if shldErr != nil {
		log.Fatalf("Error AutoMigrate ShieldApp %v", shldErr)
	}
	usrErr := db.AutoMigrate(&models.User{})
	if usrErr != nil {
		log.Fatalf("Error AutoMigrate User %v", usrErr)
	}

	ShieldAppDomainMappingErr := db.AutoMigrate(&models.ShieldAppDomainMapping{})
	if ShieldAppDomainMappingErr != nil {
		log.Fatalf("Error AutoMigrate ShieldAppDomainMapping %v", ShieldAppDomainMappingErr)
	}

	polGrpPolErr := db.AutoMigrate(&models.PolGpPolicy{})
	if polGrpPolErr != nil {
		log.Fatalf("Error  AutoMigrate PolGpPolicy %v", polGrpPolErr)
	}

	acPolGrpErr := db.AutoMigrate(&models.AcPolGrp{})
	if acPolGrpErr != nil {
		log.Fatalf("Error  AutoMigrate AcPolGrp %v", acPolGrpErr)
	}

	acPolErr := db.AutoMigrate(&models.AcPolicy{})
	if acPolErr != nil {
		log.Fatalf("Error  AutoMigrate AcPolicy %v", acPolErr)
	}

	acGpAcErr := db.AutoMigrate(&models.ActGpAction{})
	if acGpAcErr != nil {
		log.Fatalf("Error  AutoMigrate ActGpAction %v", acGpAcErr)
	}

	acAcGrpErr := db.AutoMigrate(&models.AcActGrp{})
	if acAcGrpErr != nil {
		log.Fatalf("Error  AutoMigrate AcActGrp %v", acAcGrpErr)
	}

	acAcErr := db.AutoMigrate(&models.AcAction{})
	if acAcErr != nil {
		log.Fatalf("Error  AutoMigrate AcAction %v", acAcErr)
	}

	acReAcErr := db.AutoMigrate(&models.AcResAction{})
	if acReAcErr != nil {
		log.Fatalf("Error  AutoMigrate AcResAction %v", acReAcErr)
	}

	acReGrpResErr := db.AutoMigrate(&models.AcResGpRes{})
	if acReGrpResErr != nil {
		log.Fatalf("Error  AutoMigrate AcResGpRes %v", acReGrpResErr)
	}

	acReGrpErr := db.AutoMigrate(&models.AcResGrp{})
	if acReGrpErr != nil {
		log.Fatalf("Error  AutoMigrate AcResGrp %v", acReGrpErr)
	}

	acRErr := db.AutoMigrate(&models.AcResource{})
	if acRErr != nil {
		log.Fatalf("Error  AutoMigrate AcResource %v", acRErr)
	}

	tmMbErr := db.AutoMigrate(&models.TeamMember{})
	if tmMbErr != nil {
		log.Fatalf("Error  AutoMigrate TeamMember %v", tmMbErr)
	}

	teamErr := db.AutoMigrate(&models.Team{})
	if teamErr != nil {
		log.Fatalf("Error  AutoMigrate Team %v", teamErr)
	}

	deUsrSpErr := db.AutoMigrate(&models.DefaultUserSpace{})
	if deUsrSpErr != nil {
		log.Fatalf("Error  AutoMigrate DefaultUserSpace %v", deUsrSpErr)
	}

	mbrRErr := db.AutoMigrate(&models.MemberRole{})
	if mbrRErr != nil {
		log.Fatalf("Error  AutoMigrate MemberRole  %v", mbrRErr)
	}

	rErr := db.AutoMigrate(&models.Role{})
	if rErr != nil {
		log.Fatalf("Error AutoMigrate Role %v", rErr)
	}

	acPerErr := db.AutoMigrate(&models.AcPermissions{})
	if acPerErr != nil {
		log.Fatalf("Error AutoMigrate AcPermissions %v", acPerErr)
	}
	perPolGrpErr := db.AutoMigrate(&models.PerPolGrps{})
	if perPolGrpErr != nil {
		log.Fatalf("Error AutoMigrate PerPolGrps %v", perPolGrpErr)
	}

	polGrpSubErr := db.AutoMigrate(&models.AcPolGrpSub{})
	if polGrpSubErr != nil {
		log.Fatalf("Error AutoMigrate AcPolGrpSub %v", polGrpSubErr)
	}

	InviteDetailsErr := db.AutoMigrate(&models.InviteDetails{})
	if InviteDetailsErr != nil {
		log.Fatalf("Error AutoMigrate InviteDetails %v", InviteDetailsErr)
	}

	InviteErr := db.AutoMigrate(&models.Invites{})
	if InviteErr != nil {
		log.Fatalf("Error AutoMigrate Invites %v", InviteErr)
	}

	SpaceMemberErr := db.AutoMigrate(&models.SpaceMember{})
	if SpaceMemberErr != nil {
		log.Fatalf("Error AutoMigrate SpaceMember %v", SpaceMemberErr)
	}

	EntitiesErr := db.AutoMigrate(&models.Entities{})
	if EntitiesErr != nil {
		log.Fatalf("Error AutoMigrate Entities %v", EntitiesErr)
	}

	PolGrpSubsEntityMappingErr := db.AutoMigrate(&models.PolGrpSubsEntityMapping{})
	if PolGrpSubsEntityMappingErr != nil {
		log.Fatalf("Error AutoMigrate PolGrpSubsEntityMapping %v", PolGrpSubsEntityMappingErr)
	}

	EntityTypeDefinitionErr := db.AutoMigrate(&models.EntityTypeDefinition{})
	if EntityTypeDefinitionErr != nil {
		log.Fatalf("Error AutoMigrate EntityDefinition %v", EntityTypeDefinitionErr)
	}

	EntitySpaceMappingErr := db.AutoMigrate(&models.EntitySpaceMapping{})
	if EntitySpaceMappingErr != nil {
		log.Fatalf("Error AutoMigrate EntityDefinition %v", EntitySpaceMappingErr)
	}

	//seeding default app for which app needs to be managed using shield (use the same client_id in the login request)
	defaultPermissionsErr := db.Exec(`
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('435e7c65-1fd7-4718-9944-69e90b520542','Name',NULL,NULL,'False','2022-02-11 13:57:06.069568+00','2022-02-11 13:57:06.069568+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('5e417e6d-b24c-43e3-841f-d4f5ead70ff1','Address',NULL,NULL,'False','2022-02-11 13:57:06.838323+00','2022-02-11 13:57:06.838323+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('4c6e1190-6049-4435-9292-5f58037b02fb','Phone',NULL,NULL,'False','2022-02-11 13:57:07.612243+00','2022-02-11 13:57:07.612243+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('5194c89e-2f68-499f-a521-cb086f706ba8','Calendar','Access event details from your calendar app.',NULL,'False','2022-02-11 13:57:08.37531+00','2022-02-11 13:57:08.37531+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('2feea3e9-2967-470a-892b-a578dea752d2','Location','Access your location while you are using the app.',NULL,'False','2022-02-11 13:57:09.139656+00','2022-02-11 13:57:09.139656+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('0e9e649f-94e1-4c5c-b651-3a8729d97c15','WiFi Connection Info','Full network access.',NULL,'False','2022-02-11 13:57:09.91085+00','2022-02-11 13:57:09.91085+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('fb69c178-e739-4c8a-a10b-a3ad490c23ae','Device ID',NULL,NULL,'False','2022-02-11 13:57:10.677277+00','2022-02-11 13:57:10.677277+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('dfbd8da1-0eb5-4082-8bf9-bb676450758d','Can create and delete your appblox activity.',NULL,NULL,'False','2022-02-11 13:57:11.44848+00','2022-02-11 13:57:11.44848+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('346a4734-42bb-4882-953a-3b3356d15248','Contacts',NULL,NULL,'False','2022-02-11 13:57:12.219896+00','2022-02-11 13:57:12.219896+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('34f336bb-0b71-40bf-bbc8-fa9d9965a16c','Files',NULL,NULL,'False','2022-02-11 13:57:12.991628+00','2022-02-11 13:57:12.991628+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('18d7bbbc-cc0e-4a99-a8d3-616eeff37344','Microphone',NULL,NULL,'False','2022-02-11 13:57:13.753364+00','2022-02-11 13:57:13.753364+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('960e93c5-2e2f-4a23-80ab-dff7127bce2d','Email','Read. Compose. Send your emails.',NULL,'True','2022-02-11 13:57:04.262583+00','2022-02-11 13:57:04.262583+00') on conflict do nothing;
		INSERT INTO permissions("permission_id","permission_name",description,category,mandatory,"created_at","updated_at") VALUES ('ba74d677-7dfa-4e4c-80b4-b43619ab8bc5','Username',NULL,NULL,'True','2022-02-11 13:57:05.294117+00','2022-02-11 13:57:05.294117+00') on conflict do nothing ;
		`)

	if defaultPermissionsErr.Error != nil {
		log.Fatal("Error")
	}

	clientID := os.Getenv("BLOCK_ENV_URL_CLIENT_ID")

	//seeding default app for which app needs to be managed using shield (use the same client_id in the login request)
	shieldAppQuery := fmt.Sprintf(`
	INSERT INTO public.shield_apps(
		app_id, client_id, client_secret, app_name, app_sname, description, app_url, redirect_url, app_type, created_at, updated_at, deleted_at, owner_space_id, id)
		
	VALUES (nanoid(),?, 'your_client_secret', 'test-app', 'test-app', 'Test App', 'http://localhost:3011', '{http://localhost:3011}', 2, now(), null, null, null, null)
	ON CONFLICT DO NOTHING;
    `)

	newApp := db.Exec(shieldAppQuery, clientID)
	if newApp.Error != nil {
		log.Fatal("Error")
	}

	// Query to seed  domain url mapping with shield app . Change the url according to your preferences
	shieldDomainMappingsQuery := fmt.Sprintf(`
        WITH shield_app AS (
            SELECT app_id
            FROM shield_apps 
            WHERE client_id =?
        )
        INSERT INTO public.shield_app_domain_mappings (
            id,owner_app_id,url
        ) VALUES
			(nanoid(),(SELECT app_id FROM shield_app),'http://localhost:3011')ON CONFLICT DO NOTHING;`)

	shieldDomainMappings := db.Exec(shieldDomainMappingsQuery, clientID)
	if shieldDomainMappings.Error != nil {
		log.Fatal("Error")
	}

	//seeding permission for  shield app
	newAppPermissionErr := db.Exec(`
		INSERT INTO public.app_permissions(
			app_id, permission_id, mandatory, created_at, updated_at)
		select a.app_id,p.permission_id,p.mandatory,now(),null from shield_apps a inner join permissions p on true 
		where a.client_id in (?) on conflict do nothing;
		`, clientID)

	if newAppPermissionErr.Error != nil {
		log.Fatal("Error")
	}

	//initial entity defenition seeding
	newEntityDefinitionErr := db.Exec(`
		INSERT INTO public.entity_type_definitions(id,name,display_name) VALUES(1,'entity_type_1','entity_type_1'),(2,'entity_type_2','entity_type_2'),(3,'entity_type_3','entity_type__3') on conflict do nothing;
		`)

	if newEntityDefinitionErr.Error != nil {
		log.Fatal("Error")
	}

	// Space Access Entities insertion for each entity type
	newSpaceAccessEntitiesErr := db.Exec(`
	INSERT INTO public.entities(
		entity_id, created_at, updated_at, deleted_at, type, label)
		VALUES ('1', now(), now(), null, 1, 'Space Access'),('2', now(), now(), null, 2, 'Space Access'),('3', now(), now(), null, 3, 'Space Access') on conflict do nothing
	`)

	if newSpaceAccessEntitiesErr.Error != nil {
		log.Fatal("Error")
	}

	// Query to seed resources
	acResQuery := fmt.Sprintf(`
	   
	    INSERT INTO public.ac_resources (
	        id, created_at, updated_at, deleted_at,
	        name, description, path, function_name,
	        entity_name, function_method, version, opt_counter,
	         is_authorised, is_authenticated
	    ) VALUES
			('550e8400-e29b-41d4-a716-446655440000', '2024-01-16 09:46:31.875919+00', '2024-01-16 09:46:31.875919+00', NULL, 'teams-add-permissions', NULL, '/api/spaces/v0.1/teams-add-permissions/invoke', 'list-licenses', 'spaces', 'POST', 'V.01', NULL, 2, 2) ,
			('fK8cav2w8b7mCyvoP3Bdr', '2024-01-18 06:22:38.981638+00', '2024-01-18 06:22:38.981638+00', NULL, 'Resource for list-spaces', NULL, '/api/spaces/v0.1/list-spaces/invoke', 'list-spaces', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('550e84s00-e29b-4hgy4-a716-446655440000', '2024-01-18 10:35:49.613494+00', '2024-01-18 10:35:49.613494+00', NULL, 'create-invite_link', NULL, '/api/spaces/v0.1/create-invite-link/invoke', 'create-invite_link', 'spaces', 'POST', 'V.01', NULL, 2, 2),
			('cYFau9EH6ANQZC-LrvCHO', '2024-01-18 11:08:51.482964+00', '2024-01-18 11:08:51.482964+00', NULL, 'Resource for teams-list-to-add-pol-grp-subs', NULL, '/api/spaces/v0.1/teams-list-to-add-pol-grp-subs/invoke', 'teams-list-to-add-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('FNTFcHkDGRYrU9-m_gkSA', '2024-01-18 11:08:51.490274+00', '2024-01-18 11:08:51.490274+00', NULL, 'Resource for user-list-available-entities', NULL, '/api/spaces/v0.1/user-list-available-entities/invoke', 'user-list-available-entities', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('jzAxzMT81Kv8aCZEcQgc8', '2024-01-18 11:08:51.491261+00', '2024-01-18 11:08:51.491261+00', NULL, 'Resource for check-role-name', NULL, '/api/spaces/v0.1/check-role-name/invoke', 'check-role-name', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('9LjrDvq1AFKTcv8hBbUGP', '2024-01-18 11:08:51.492943+00', '2024-01-18 11:08:51.492943+00', NULL, 'Resource for cancel-invite', NULL, '/api/spaces/v0.1/cancel-invite/invoke', 'cancel-invite', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('d-wOysylDamL7hIO0QrZC', '2024-01-18 11:08:51.496833+00', '2024-01-18 11:08:51.496833+00', NULL, 'Resource for roles-list-existing-pol-grp-subs', NULL, '/api/spaces/v0.1/roles-list-existing-pol-grp-subs/invoke', 'roles-list-existing-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('h0Iy61Gqa8AOwq1gt2RUx', '2024-01-18 11:08:51.498867+00', '2024-01-18 11:08:51.498867+00', NULL, 'Resource for roles-delete-existing-pol-grp-subs', NULL, '/api/spaces/v0.1/roles-delete-existing-pol-grp-subs/invoke', 'roles-delete-existing-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('84_Yq5VzvLiVEBmiCfsh2', '2024-01-18 11:08:51.501083+00', '2024-01-18 11:08:51.501083+00', NULL, 'Resource for teams-add-pol-grp-subs', NULL, '/api/spaces/v0.1/teams-add-pol-grp-subs/invoke', 'teams-add-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('DA6eNzlrdKH5EqjpYjjga', '2024-01-18 11:08:51.504112+00', '2024-01-18 11:08:51.504112+00', NULL, 'Resource for teams-list-available-permissions', NULL, '/api/spaces/v0.1/teams-list-available-permissions/invoke', 'teams-list-available-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('N4sjxTnNf1hewEXqWXWJL', '2024-01-18 11:08:51.507645+00', '2024-01-18 11:08:51.507645+00', NULL, 'Resource for teams-add-permissions', NULL, '/api/spaces/v0.1/teams-add-permissions/invoke', 'teams-add-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('xWLAvUnRa6YsH9ueNnc9_', '2024-01-18 11:08:51.510679+00', '2024-01-18 11:08:51.510679+00', NULL, 'Resource for create-logo-signed-url', NULL, '/api/spaces/v0.1/create-logo-signed-url/invoke', 'create-logo-signed-url', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('KrzwOoRAQ-Gb0nV4v7X7R', '2024-01-18 11:08:51.515918+00', '2024-01-18 11:08:51.515918+00', NULL, 'Resource for roles-search-user', NULL, '/api/spaces/v0.1/roles-search-user/invoke', 'roles-search-user', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('QMo-wqpa6ptG9Uz-JphVt', '2024-01-18 11:08:51.517827+00', '2024-01-18 11:08:51.517827+00', NULL, 'Resource for list-spaces-detailed', NULL, '/api/spaces/v0.1/list-spaces-detailed/invoke', 'list-spaces-detailed', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('ybT6UL9cWQOvcUBfSIFRj', '2024-01-18 11:08:51.519145+00', '2024-01-18 11:08:51.519145+00', NULL, 'Resource for send-user-invite-email', NULL, '/api/spaces/v0.1/send-user-invite-email/invoke', 'send-user-invite-email', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('gmFwBdvYaUY_BvbHMNhyX', '2024-01-18 11:08:51.521206+00', '2024-01-18 11:08:51.521206+00', NULL, 'Resource for list-invited-users', NULL, '/api/spaces/v0.1/list-invited-users/invoke', 'list-invited-users', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('g5fh9UFlrq_1OP1IgZksH', '2024-01-18 11:08:51.525172+00', '2024-01-18 11:08:51.525172+00', NULL, 'Resource for user-list-permissions', NULL, '/api/spaces/v0.1/user-list-permissions/invoke', 'user-list-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('rKsvAsb9hAeN4vNoqdAnM', '2024-01-18 11:08:51.531688+00', '2024-01-18 11:08:51.531688+00', NULL, 'Resource for teams-list-permissions', NULL, '/api/spaces/v0.1/teams-list-permissions/invoke', 'teams-list-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('fQKcNySVMXeqbIFNx8oMW', '2024-01-18 11:08:51.532744+00', '2024-01-18 11:08:51.532744+00', NULL, 'Resource for list-users', NULL, '/api/spaces/v0.1/list-users/invoke', 'list-users', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('SjSs4Hvbya78GUJUTc3M0', '2024-01-18 11:08:51.534711+00', '2024-01-18 11:08:51.534711+00', NULL, 'Resource for user-list-pol-grp-subs-from-roles', NULL, '/api/spaces/v0.1/user-list-pol-grp-subs-from-roles/invoke', 'user-list-pol-grp-subs-from-roles', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('66I5V-LFfq9YoUtYUewBG', '2024-01-18 11:08:51.536351+00', '2024-01-18 11:08:51.536351+00', NULL, 'Resource for user-list-to-add-pol-grp-subs', NULL, '/api/spaces/v0.1/user-list-to-add-pol-grp-subs/invoke', 'user-list-to-add-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('5J5NUEXeWR1ASgog2ZO59', '2024-01-18 11:08:51.537735+00', '2024-01-18 11:08:51.537735+00', NULL, 'Resource for set-inuse-block-in-app', NULL, '/api/spaces/v0.1/set-inuse-block-in-app/invoke', 'set-inuse-block-in-app', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('URAEQKuj1bzIzsnOIQm5z', '2024-01-18 11:08:51.540554+00', '2024-01-18 11:08:51.540554+00', NULL, 'Resource for update-assign-block-to-app', NULL, '/api/spaces/v0.1/update-assign-block-to-app/invoke', 'update-assign-block-to-app', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('wx9GgyMP9H51e8p69CxDo', '2024-01-18 11:08:51.54207+00', '2024-01-18 11:08:51.54207+00', NULL, 'Resource for roles-add-permissions', NULL, '/api/spaces/v0.1/roles-add-permissions/invoke', 'roles-add-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('gbpiyFlalU_l_YxRmDuo6', '2024-01-18 11:08:51.544261+00', '2024-01-18 11:08:51.544261+00', NULL, 'Resource for create-space', NULL, '/api/spaces/v0.1/create-space/invoke', 'create-space', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('GGI_3qpFjzPjl3rExq6Mk', '2024-01-18 11:08:51.546857+00', '2024-01-18 11:08:51.546857+00', NULL, 'Resource for check-business-name', NULL, '/api/spaces/v0.1/check-business-name/invoke', 'check-business-name', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('ThEYW6j1K4pIatwyRRoYv', '2024-01-18 11:08:51.55003+00', '2024-01-18 11:08:51.55003+00', NULL, 'Resource for revoke-invite', NULL, '/api/spaces/v0.1/revoke-invite/invoke', 'revoke-invite', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('joDRJArzH6tcyai6sfale', '2024-01-18 11:08:51.552815+00', '2024-01-18 11:08:51.552815+00', NULL, 'Resource for teams-list-existing-pol-grp-subs', NULL, '/api/spaces/v0.1/teams-list-existing-pol-grp-subs/invoke', 'teams-list-existing-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('fERUBEMH58hf03aW7ERCK', '2024-01-18 11:08:51.554045+00', '2024-01-18 11:08:51.554045+00', NULL, 'Resource for user-list-pol-grp-subs-from-teams', NULL, '/api/spaces/v0.1/user-list-pol-grp-subs-from-teams/invoke', 'user-list-pol-grp-subs-from-teams', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('-KIBHqjXbpmFmZ-62cgBg', '2024-01-18 11:08:51.559187+00', '2024-01-18 11:08:51.559187+00', NULL, 'Resource for user-add-permissions', NULL, '/api/spaces/v0.1/user-add-permissions/invoke', 'user-add-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('Hp0MOQJYJOKhxlRTgG2tN', '2024-01-18 11:08:51.561675+00', '2024-01-18 11:08:51.561675+00', NULL, 'Resource for roles-list-permissions', NULL, '/api/spaces/v0.1/roles-list-permissions/invoke', 'roles-list-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('X6s_JM330WCXad9tCxmkm', '2024-01-18 11:08:51.565457+00', '2024-01-18 11:08:51.565457+00', NULL, 'Resource for roles-list-users', NULL, '/api/spaces/v0.1/roles-list-users/invoke', 'roles-list-users', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('1E40HEwoY0ltpt3A6qEXD', '2024-01-18 11:08:51.566302+00', '2024-01-18 11:08:51.566302+00', NULL, 'Resource for create-role', NULL, '/api/spaces/v0.1/create-role/invoke', 'create-role', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('ZxV7YrY4b33-nZz-ECSZt', '2024-01-18 11:08:51.568668+00', '2024-01-18 11:08:51.568668+00', NULL, 'Resource for update-user', NULL, '/api/spaces/v0.1/update-user/invoke', 'update-user', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('n-NbLgT86s2Agrd6Ba-Zc', '2024-01-18 11:08:51.570122+00', '2024-01-18 11:08:51.570122+00', NULL, 'Resource for teams-delete-user', NULL, '/api/spaces/v0.1/teams-delete-user/invoke', 'teams-delete-user', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('G7UPVuJUHd36plGqSSJNM', '2024-01-18 11:08:51.57299+00', '2024-01-18 11:08:51.57299+00', NULL, 'Resource for roles-create-invite-link', NULL, '/api/spaces/v0.1/roles-create-invite-link/invoke', 'roles-create-invite-link', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('0xLmZz6GSOulH_ttxrxoS', '2024-01-18 11:08:51.576255+00', '2024-01-18 11:08:51.576255+00', NULL, 'Resource for unassign-block-from-app', NULL, '/api/spaces/v0.1/unassign-block-from-app/invoke', 'unassign-block-from-app', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('Z_k8T1TiLY-5Jwj-aOPJu', '2024-01-18 11:08:51.577691+00', '2024-01-18 11:08:51.577691+00', NULL, 'Resource for list-teams', NULL, '/api/spaces/v0.1/list-teams/invoke', 'list-teams', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('izLUFyWmgacyY5T_kWVbC', '2024-01-18 11:08:51.579879+00', '2024-01-18 11:08:51.579879+00', NULL, 'Resource for resend-invite-email', NULL, '/api/spaces/v0.1/resend-invite-email/invoke', 'resend-invite-email', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('Uj-miSHlxeBGfezqyjVkR', '2024-01-18 11:08:51.582047+00', '2024-01-18 11:08:51.582047+00', NULL, 'Resource for roles-list-entities', NULL, '/api/spaces/v0.1/roles-list-entities/invoke', 'roles-list-entities', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('VHGxvTMzg4ifDj2-Vini_', '2024-01-18 11:08:51.585959+00', '2024-01-18 11:08:51.585959+00', NULL, 'Resource for teams-search-user', NULL, '/api/spaces/v0.1/teams-search-user/invoke', 'teams-search-user', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('I-Ntfv4e48G1aQ4mkVeSb', '2024-01-18 11:08:51.586932+00', '2024-01-18 11:08:51.586932+00', NULL, 'Resource for teams-send-user-invite-email', NULL, '/api/spaces/v0.1/teams-send-user-invite-email/invoke', 'teams-send-user-invite-email', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('6_q6lzLxkWfiMA-QezCJY', '2024-01-18 11:08:51.590043+00', '2024-01-18 11:08:51.590043+00', NULL, 'Resource for create-team', NULL, '/api/spaces/v0.1/create-team/invoke', 'create-team', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('xn0BG_Bu2x_5OVxfgRhCw', '2024-01-18 11:08:51.592117+00', '2024-01-18 11:08:51.592117+00', NULL, 'Resource for delete-team', NULL, '/api/spaces/v0.1/delete-team/invoke', 'delete-team', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('TTdMbqj7iPg41viT5faok', '2024-01-18 11:08:51.595837+00', '2024-01-18 11:08:51.595837+00', NULL, 'Resource for user-list-available-permissions', NULL, '/api/spaces/v0.1/user-list-available-permissions/invoke', 'user-list-available-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('bgpsl8TkgQcP6Jif0Tiaq', '2024-01-18 11:08:51.598532+00', '2024-01-18 11:08:51.598532+00', NULL, 'Resource for user-delete-existing-pol-grp-subs', NULL, '/api/spaces/v0.1/user-delete-existing-pol-grp-subs/invoke', 'user-delete-existing-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('IiZ4BWuUEJ6GBkvJLEoEQ', '2024-01-18 11:08:51.604479+00', '2024-01-18 11:08:51.604479+00', NULL, 'Resource for list-roles', NULL, '/api/spaces/v0.1/list-roles/invoke', 'list-roles', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('DvjcCCfkYwkqLrl-S4rHE', '2024-01-18 11:08:51.611528+00', '2024-01-18 11:08:51.611528+00', NULL, 'Resource for get-user-by-id', NULL, '/api/spaces/v0.1/get-user-by-id/invoke', 'get-user-by-id', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('TxFqMz-pzhPefJXaTGuem', '2024-01-18 11:08:51.61403+00', '2024-01-18 11:08:51.61403+00', NULL, 'Resource for user-list-existing-pol-grp-subs', NULL, '/api/spaces/v0.1/user-list-existing-pol-grp-subs/invoke', 'user-list-existing-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('0lgUjlGf5jcRUrYbegzI9', '2024-01-18 11:08:51.614923+00', '2024-01-18 11:08:51.614923+00', NULL, 'Resource for user-add-pol-grp-subs', NULL, '/api/spaces/v0.1/user-add-pol-grp-subs/invoke', 'user-add-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('MKqIb6eY0ywyjS8wtFHlM', '2024-01-18 11:08:51.615693+00', '2024-01-18 11:08:51.615693+00', NULL, 'Resource for teams-create-invite-link', NULL, '/api/spaces/v0.1/teams-create-invite-link/invoke', 'teams-create-invite-link', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('WNMYOh4kuW6U71BnBFT40', '2024-01-18 11:08:51.616345+00', '2024-01-18 11:08:51.616345+00', NULL, 'Resource for teams-list-users', NULL, '/api/spaces/v0.1/teams-list-users/invoke', 'teams-list-users', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('BB0XAr6sb6_eRpL-udogo', '2024-01-18 11:08:51.616985+00', '2024-01-18 11:08:51.616985+00', NULL, 'Resource for update-role', NULL, '/api/spaces/v0.1/update-role/invoke', 'update-role', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('tlkmt2nypUKIVLrsZ2rZC', '2024-01-18 11:08:51.617835+00', '2024-01-18 11:08:51.617835+00', NULL, 'Resource for roles-list-to-add-pol-grp-subs', NULL, '/api/spaces/v0.1/roles-list-to-add-pol-grp-subs/invoke', 'roles-list-to-add-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('YWqUUZqNAKmnzgcyYOGxJ', '2024-01-18 11:08:51.619778+00', '2024-01-18 11:08:51.619778+00', NULL, 'Resource for search-user', NULL, '/api/spaces/v0.1/search-user/invoke', 'search-user', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('15lRe9db_hsVJZpat62MY', '2024-01-18 11:08:51.622376+00', '2024-01-18 11:08:51.622376+00', NULL, 'Resource for teams-delete-existing-pol-grp-subs', NULL, '/api/spaces/v0.1/teams-delete-existing-pol-grp-subs/invoke', 'teams-delete-existing-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('NyIsbOMFB-Zq8y6fikhT2', '2024-01-18 11:08:51.634359+00', '2024-01-18 11:08:51.634359+00', NULL, 'Resource for check-space-name', NULL, '/api/spaces/v0.1/check-space-name/invoke', 'check-space-name', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('QWzsvCNBM3gkQbu-wuOKA', '2024-01-18 11:08:51.63615+00', '2024-01-18 11:08:51.63615+00', NULL, 'Resource for teams-list-entities', NULL, '/api/spaces/v0.1/teams-list-entities/invoke', 'teams-list-entities', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('wswn9DbRFAhv3FtIU7I9I', '2024-01-18 11:08:51.63869+00', '2024-01-18 11:08:51.63869+00', NULL, 'Resource for user-list-entities', NULL, '/api/spaces/v0.1/user-list-entities/invoke', 'user-list-entities', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('J3voHDKaM-SfjwH7B5zHZ', '2024-01-18 11:08:51.648064+00', '2024-01-18 11:08:51.648064+00', NULL, 'Resource for check-assigned-block-to-app', NULL, '/api/spaces/v0.1/check-assigned-block-to-app/invoke', 'check-assigned-block-to-app', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('UfkuayXNbQq41etTK2qdQ', '2024-01-18 11:08:51.650737+00', '2024-01-18 11:08:51.650737+00', NULL, 'Resource for update-space', NULL, '/api/spaces/v0.1/update-space/invoke', 'update-space', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('Fii5tyAdtHMLEGxNTpi9k', '2024-01-18 11:08:51.656374+00', '2024-01-18 11:08:51.656374+00', NULL, 'Resource for create-invite-link', NULL, '/api/spaces/v0.1/create-invite-link/invoke', 'create-invite-link', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('1sPvavaQygX4Mx8MkmSLH', '2024-01-18 11:08:51.660287+00', '2024-01-18 11:08:51.660287+00', NULL, 'Resource for roles-add-pol-grp-subs', NULL, '/api/spaces/v0.1/roles-add-pol-grp-subs/invoke', 'roles-add-pol-grp-subs', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('udbHDGgGsvrp3x8C45Ewm', '2024-01-18 11:08:51.662441+00', '2024-01-18 11:08:51.662441+00', NULL, 'Resource for get-space-by-id', NULL, '/api/spaces/v0.1/get-space-by-id/invoke', 'get-space-by-id', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('jqVlE5nwv5SsDBssJFbkb', '2024-01-18 11:08:51.66505+00', '2024-01-18 11:08:51.66505+00', NULL, 'Resource for roles-send-user-invite-email', NULL, '/api/spaces/v0.1/roles-send-user-invite-email/invoke', 'roles-send-user-invite-email', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('ZOAwbl-C9xwew7QJwG7dG', '2024-01-18 11:08:51.666435+00', '2024-01-18 11:08:51.666435+00', NULL, 'Resource for check-team-name', NULL, '/api/spaces/v0.1/check-team-name/invoke', 'check-team-name', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('CHznnsnt5Bm2163PFY5fz', '2024-01-18 11:08:51.674641+00', '2024-01-18 11:08:51.674641+00', NULL, 'Resource for assign-block-to-app', NULL, '/api/spaces/v0.1/assign-block-to-app/invoke', 'assign-block-to-app', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('Si25JQ8ddYkURkv_CSrlw', '2024-01-18 11:08:51.682028+00', '2024-01-18 11:08:51.682028+00', NULL, 'Resource for roles-list-available-permissions', NULL, '/api/spaces/v0.1/roles-list-available-permissions/invoke', 'roles-list-available-permissions', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('LJGas2LZHBNTcdm_W_GTL', '2024-01-18 11:08:51.68639+00', '2024-01-18 11:08:51.68639+00', NULL, 'Resource for delete-space', NULL, '/api/spaces/v0.1/delete-space/invoke', 'delete-space', 'spaces', 'POST', 'v0.1', NULL , 2, 2),
			('9eR07Sa5utsH1NXov9TZt', '2024-01-18 11:08:51.690547+00', '2024-01-18 11:08:51.690547+00', NULL, 'Resource for update-team', NULL, '/api/spaces/v0.1/update-team/invoke', 'update-team', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('xq1B9es0X7WxpiVYN9csM', '2024-01-18 11:08:51.699858+00', '2024-01-18 11:08:51.699858+00', NULL, 'Resource for get-invite-by-id', NULL, '/api/spaces/v0.1/get-invite-by-id/invoke', 'get-invite-by-id', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('LtXPeRYgWGibsPZlowWe2', '2024-01-18 11:08:51.701927+00', '2024-01-18 11:08:51.701927+00', NULL, 'Resource for roles-delete-user', NULL, '/api/spaces/v0.1/roles-delete-user/invoke', 'roles-delete-user', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('tWQHiQJrI0U190O-hd7sj', '2024-01-18 11:08:51.654099+00', '2024-01-18 11:08:51.654099+00', NULL, 'Resource for accept-invite', NULL, '/api/spaces/v0.1/accept-invite/invoke', 'accept-invite', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('twerfHiyudjkI0U190O-hd7sj', '2024-01-18 11:08:51.664079+00', '2024-01-18 11:08:51.653069+00', NULL, 'Resource for list-entity-definition', NULL, '/api/spaces/v0.1/list-entity-definition/invoke', 'list-entity-definition', 'spaces', 'POST', 'v0.1', NULL, 1, 2),
			('QNXip-DHCwaO6TP1YbuU', '2024-01-18 11:08:51.664079+00', '2024-01-18 11:08:51.653069+00', NULL, 'Resource for create-entity', NULL, '/api/spaces/v0.1/create-entity/invoke', 'create-entity', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('QNXibdtb-DHCwaO6TP1YbuU', '2024-01-18 11:08:51.664079+00', '2024-01-18 11:08:51.653069+00', NULL, 'Resource for user-add-entities', NULL, '/api/spaces/v0.1/user-add-entities/invoke', 'user-add-entities', 'spaces', 'POST', 'v0.1', NULL, 2, 2),
			('QNXibdtb-DygdwaO6TP1YbuU', '2024-01-18 11:08:51.664079+00', '2024-01-18 11:08:51.653069+00', NULL, 'Resource for entity name check', NULL, '/api/spaces/v0.1/check-entity-name/invoke', 'check-entity-name', 'spaces', 'POST', 'v0.1', NULL, 2, 2)
	        ON CONFLICT DO NOTHING;`)

	acRes := db.Exec(acResQuery)
	if acRes.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_res_grps  - we are grouping our resources using secondary table called ac_res_grps
	acResGrps := db.Exec(`INSERT INTO public.ac_res_grps(
	id, created_at, updated_at, deleted_at, owner_space_id, name, description, is_predefined, opt_counter, type)
	VALUES ('FrAbAt75VZf8vZjKPVzh-',now(),now(),null,null,'ACL Access Permission','Resource group for ACL Access Permission', true,null,1) on conflict do nothing;`)

	if acResGrps.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_res_gp_res - we have a bridge table for ac_resources and ac_res_grps
	acResGrpRes := db.Exec(`INSERT INTO public.ac_res_gp_res(
		id, created_at, updated_at, deleted_at, ac_res_grp_id, ac_resource_id, opt_counter)
		select nanoid(),now(),now(),null,ac_res_grp.id,ac_res.id,null
		from ac_res_grps ac_res_grp left join ac_resources ac_res on true
		where ac_res_grp.name in ('ACL Access Permission') and
		ac_res.function_name in ('get-user-details','list-spaces','get-space-by-id','list-users','list-roles','list-teams','user-list-permissions',
		'user-list-existing-pol-grp-subs','user-list-to-add-pol-grp-subs','user-add-pol-grp-subs','user-list-available-permissions',
		'user-delete-existing-pol-grp-subs','user-list-app-entities','user-add-app-entities','user-list-entities','user-add-block-entities','create-invite-link','roles-search-user','send-user-invite-email',
		'roles-list-users','roles-list-existing-pol-grp-subs','roles-list-to-add-pol-grp-subs','roles-add-pol-grp-subs','check-role-name','check-team-name','check-entity-name'
		,'roles-delete-existing-pol-grp-subs','roles-list-app-entities','roles-add-app-entities','roles-list-entities','roles-add-block-entities','roles-create-invite-link','roles-send-user-invite-email',
		'teams-list-users','teams-create-invite-link','teams-search-user','teams-send-user-invite-email','teams-list-existing-pol-grp-subs','teams-list-to-add-pol-grp-subs','teams-add-pol-grp-subs',
		'teams-delete-existing-pol-grp-subs','teams-list-app-entities','user-add-permissions','user-add-entities',
		'teams-add-app-entities','teams-list-entities','teams-add-block-entities','create-role','create-team','get-user-by-id','update-user','list-spaces-detailed','roles-delete-user','cancel_invite','list-entity-definition','user-list-available-entities','teams-list-available-permissions','teams-add-permissions','teams-list-permissions','roles-list-available-permissions','roles-add-permissions','roles-list-permissions')  on conflict do nothing;`)

	if acResGrpRes.Error != nil {
		log.Fatal("Error")
	}
	//seeding actions table its basically api action methods
	acAct := db.Exec(`
	
	INSERT INTO public.ac_actions(
			id, created_at, updated_at, deleted_at,  name, description,opt_counter)
			VALUES ('QynaTw021PRwJg57GauYG',now(),now(),null,'invoke','predefined invoke action for resource',null) on conflict do nothing;`)

	if acAct.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_act_grps  secondary table
	acActGrps := db.Exec(`INSERT INTO public.ac_act_grps(
	   id, created_at, updated_at, deleted_at, owner_space_id, description, name, is_predefined, opt_counter, type)
	   VALUES ('KXjyu1O5oCNS00t9oMX-y',now(),now(),null,null,'Action group for ACL Access Permission','ACL Access Permission',true,null,1)on conflict do nothing;`)

	if acActGrps.Error != nil {
		log.Fatal("Error")
	}
	//seeding act_gp_actions and ac_actions bridge table
	actGrpAct := db.Exec(`INSERT INTO public.act_gp_actions(
		id, created_at, updated_at, deleted_at, ac_act_grp_id, ac_action_id, opt_counter)
		select 'EPATshPgmCI2PZC9Oec9_',now(),now(),null,actgrp.id,ac_act.id,null from ac_act_grps actgrp
		left join ac_actions ac_act on true
		where actgrp.name in ('ACL Access Permission')
	    and ac_act.description in ('predefined invoke action for resource')  on conflict do nothing;`)

	if actGrpAct.Error != nil {
		log.Fatal("Error")
	}

	// seeding ac_policies - we are grouping multiple resources to policy
	acPol := db.Exec(`INSERT INTO public.ac_policies(
		id, created_at, updated_at, deleted_at, ac_act_grp_id, ac_res_grp_id, owner_space_id, created_by, updated_by, name, description, path, opt_counter, is_predefined, type)
	    select 'ErKryecbBKilzx2K-CHv1',now(),now(),null,actgrp.id,acresgrp.id,null,null,null,'ACL Access Permission','Ac Policy for ACL Access Permission',null,null,true,1
	    from ac_act_grps actgrp left join ac_res_grps acresgrp on true
	    where actgrp.name in ('ACL Access Permission') and acresgrp.name in ('ACL Access Permission') on conflict do nothing;`)

	if acPol.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_pol_grps  table - secondary table for ac_policies
	actPolGrp := db.Exec(`INSERT INTO public.ac_pol_grps(
		id, created_at, updated_at, deleted_at, owner_space_id, description, name, opt_counter, is_predefined, type, entity_type, display_name, entity_types)
	    select 'sUm4zZvsT1u-H6A_d4caH',now(),now(),null,null,'Ac Policy Group for ACL Access','ACL Access',null,true,1,null,'ACL Access','{1,2,3}' on conflict do nothing;`)

	if actPolGrp.Error != nil {
		log.Fatal("Error")
	}

	//seeding pol_gp_policies - ac_pol_grps and ac_policies bridge table
	polGrpPol := db.Exec(`INSERT INTO public.pol_gp_policies(
		id, created_at, updated_at, deleted_at, ac_policy_id, ac_pol_grp_id, opt_counter)
		select '0LnzPpsybC5Fl2jnlfmEy',now(),now(),null,acpol.id,acpolgrp.id,null
		from ac_policies acpol left join ac_pol_grps acpolgrp on true
		where acpol.name in ('ACL Access Permission') and acpolgrp.name in ('ACL Access') on conflict do nothing;`)

	if polGrpPol.Error != nil {
		log.Fatal("Error")
	}

	//seeding permissions table - permissions basically API usage access
	acPer := db.Exec(`INSERT INTO public.ac_permissions(
		id, created_at, updated_at, deleted_at, description, name, opt_counter, is_predefined, type, display_name)
	    select 'IOE6rugUKHCsbtx655ZDi',now(),now(),null,'Ac Permission for ACL Access Permission','ACL Access Permission',null,true,1,'ACL Access Permission'
	    on conflict do nothing;`)

	if acPer.Error != nil {
		log.Fatal("Error")
	}

	//seeding per_pol_grps - ac_permissions and  policies connection table
	perPolGrp := db.Exec(`INSERT INTO public.per_pol_grps(
	    id, created_at, updated_at, deleted_at, ac_permission_id, ac_pol_grp_id)
	    select '3JQaaRWL71tHDtdgPbiyE',now(),now(),null,acper.id,acpolgrp.id from ac_permissions acper
	    left join ac_pol_grps acpolgrp on true
	    where acper.name in ('ACL Access Permission') and acpolgrp.name in ('ACL Access')
		on conflict do nothing;`)

	if perPolGrp.Error != nil {
		log.Fatal("Error")
	}

	// Registering a sample todo app and adding resources and creating permissions for granting resources

	todoAppClientID := os.Getenv("TODO_APP_CLIENT_ID")

	//seeding default app for which app needs to be managed using shield (use the same client_id in the login request)
	todoAppRegisterQuery := fmt.Sprintf(`
	INSERT INTO public.shield_apps(
		app_id, client_id, client_secret, app_name, app_sname, description, app_url, redirect_url, app_type, created_at, updated_at, deleted_at, owner_space_id, id)
		
	VALUES (nanoid(),?, '^\uM8a+£hUgCj=N_krV>0?:qI[K0p-qCl-1Upe', 'todo-app', 'todo-app', 'Test App', 'http://localhost:3001', '{http://localhost:3001}', 2, now(), null, null, null, null)
	ON CONFLICT DO NOTHING;
    `)

	todoApp := db.Exec(todoAppRegisterQuery, todoAppClientID)
	if todoApp.Error != nil {
		log.Fatal("Error")
	}

	// Query to seed  domain url mapping with shield app . Change the url according to your preferences
	todoAppDomainMapingQuery := fmt.Sprintf(`
        WITH shield_app AS (
            SELECT app_id
            FROM shield_apps 
            WHERE client_id =?
        )
        INSERT INTO public.shield_app_domain_mappings (
            id,owner_app_id,url
        ) VALUES
			(nanoid(),(SELECT app_id FROM shield_app),'http://localhost:3001')ON CONFLICT DO NOTHING;`)

	todoAppDomainMappings := db.Exec(todoAppDomainMapingQuery, todoAppClientID)
	if todoAppDomainMappings.Error != nil {
		log.Fatal("Error")
	}

	//seeding permission for  shield app
	todoAppPermissionErr := db.Exec(`
		INSERT INTO public.app_permissions(
			app_id, permission_id, mandatory, created_at, updated_at)
		select a.app_id,p.permission_id,p.mandatory,now(),null from shield_apps a inner join permissions p on true 
		where a.client_id in (?) on conflict do nothing;
		`, todoAppClientID)

	if todoAppPermissionErr.Error != nil {
		log.Fatal("Error")
	}

	//initial entity defenition seeding

	// Query to seed resources
	todoAcResQuery := fmt.Sprintf(`
	    WITH shield_app AS (
	        SELECT app_id
	        FROM shield_apps
	        WHERE client_id =?
	    )
	    INSERT INTO public.ac_resources (
	        id, created_at, updated_at, deleted_at,
	        name, description, path, function_name,
	        entity_name, function_method, version, opt_counter,
	         is_authorised, is_authenticated
	    ) VALUES
			('6qpQOAUO5I8rIf8_0N83vONLsRVUsoJ-', '2024-01-16 09:46:31.875919+00', '2024-01-16 09:46:31.875919+00', NULL, 'addToDo', NULL, '/api/todo/v0.1/addToDo/invoke', 'addToDo', 'todo', 'POST', 'V.01', NULL, 2, 2) ,('2jKGc8nZdRagQTnRdM1tVsujjL6j6DTz', '2024-01-16 09:46:31.875919+00', '2024-01-16 09:46:31.875919+00', NULL, 'listToDo', NULL, '/api/todo/v0.1/listToDo/invoke', 'listToDo', 'todo', 'POST', 'V.01', NULL, 2, 2) ,('Vl1F8lTUQ7bOHsoY9flFpds6hLjdU7HO', '2024-01-16 09:46:31.875919+00', '2024-01-16 09:46:31.875919+00', NULL, 'removeToDo', NULL, '/api/todo/v0.1/removeToDo/invoke', 'removeToDo', 'todo', 'POST', 'V.01', NULL, 2, 2) 
			
	        ON CONFLICT DO NOTHING;`)

	todoRes := db.Exec(todoAcResQuery, todoAppClientID)
	if todoRes.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_res_grps  - we are grouping our resources using secondary table called ac_res_grps
	todoResGrps := db.Exec(`INSERT INTO public.ac_res_grps(
	id, created_at, updated_at, deleted_at, owner_space_id, name, description, is_predefined, opt_counter, type)
	VALUES ('SHBhPKd1i2-mWX_o9B4udR9pS0VGL8i0',now(),now(),null,null,'ToDo List Access','Resource group for ToDo List', true,null,1),
	('D0oF6BQWomjDdJ9Bnix9WlHcIWxPjxbt',now(),now(),null,null,'ToDo Delete Access','Resource group for ToDo Delete', true,null,1),
	('NU3zbCdcascxTfa2N1mZCUTj18M_R3k0',now(),now(),null,null,'ToDo Create Access','Resource group for ToDo Create', true,null,1) on conflict do nothing;`)

	if todoResGrps.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_res_gp_res - we have a bridge table for ac_resources and ac_res_grps
	todoListResGpRes := db.Exec(`INSERT INTO public.ac_res_gp_res(
		id, created_at, updated_at, deleted_at, ac_res_grp_id, ac_resource_id, opt_counter)
		select nanoid(),now(),now(),null,ac_res_grp.id,ac_res.id,null
		from ac_res_grps ac_res_grp left join ac_resources ac_res on true
		where ac_res_grp.name in ('ToDo List Access') and
		ac_res.function_name in ('listToDo')  on conflict do nothing;`)

	if todoListResGpRes.Error != nil {
		log.Fatal("Error")
	}

	todoCreateResGpRes := db.Exec(`INSERT INTO public.ac_res_gp_res(
		id, created_at, updated_at, deleted_at, ac_res_grp_id, ac_resource_id, opt_counter)
		select nanoid(),now(),now(),null,ac_res_grp.id,ac_res.id,null
		from ac_res_grps ac_res_grp left join ac_resources ac_res on true
		where ac_res_grp.name in ('ToDo Create Access') and
		ac_res.function_name in ('addToDo')  on conflict do nothing;`)

	if todoCreateResGpRes.Error != nil {
		log.Fatal("Error")
	}

	todoDeleteResGpRes := db.Exec(`INSERT INTO public.ac_res_gp_res(
		id, created_at, updated_at, deleted_at, ac_res_grp_id, ac_resource_id, opt_counter)
		select nanoid(),now(),now(),null,ac_res_grp.id,ac_res.id,null
		from ac_res_grps ac_res_grp left join ac_resources ac_res on true
		where ac_res_grp.name in ('ToDo Delete Access') and
		ac_res.function_name in ('removeToDo')  on conflict do nothing;`)

	if todoDeleteResGpRes.Error != nil {
		log.Fatal("Error")
	}
	//seeding actions table its basically api action methods
	toDoAct := db.Exec(`

	INSERT INTO public.ac_actions(
			id, created_at, updated_at, deleted_at,  name, description,opt_counter)
			VALUES ('cssMjXvSm-sbLG1FZgH4UKU3iDQvp8rK',now(),now(),null,'invoke','predefined invoke action for resource',null) on conflict do nothing;`)

	if toDoAct.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_act_grps  secondary table
	todoActGrp := db.Exec(`INSERT INTO public.ac_act_grps(
	   id, created_at, updated_at, deleted_at, owner_space_id, description, name, is_predefined, opt_counter, type)
	   VALUES ('BZw8nyBZVhZdbwkF_izE6QMq9b5pEvXD',now(),now(),null,null,'Action group for ToDo','ToDo Action Group',true,null,1)on conflict do nothing;`)

	if todoActGrp.Error != nil {
		log.Fatal("Error")
	}
	//seeding act_gp_actions and ac_actions bridge table
	todoActGrpAct := db.Exec(`INSERT INTO public.act_gp_actions(
		id, created_at, updated_at, deleted_at, ac_act_grp_id, ac_action_id, opt_counter)
		select 'BEWYfrJtDctmQq-Nc-1k5UrvnB3UF1sG',now(),now(),null,actgrp.id,ac_act.id,null from ac_act_grps actgrp
		left join ac_actions ac_act on true
		where actgrp.name in ('ToDo Action Group')
	    and ac_act.description in ('predefined invoke action for resource')  on conflict do nothing;`)

	if todoActGrpAct.Error != nil {
		log.Fatal("Error")
	}

	// seeding ac_policies - we are grouping multiple resources to policy
	toDoListPol := db.Exec(`INSERT INTO public.ac_policies(
		id, created_at, updated_at, deleted_at, ac_act_grp_id, ac_res_grp_id, owner_space_id, created_by, updated_by, name, description, path, opt_counter, is_predefined, type)
	    select 'AfBUi3x9kEcVR2NqJoPg7eEgPOgQnNUr',now(),now(),null,actgrp.id,acresgrp.id,null,null,null,'ToDo List','Ac Policy for ToDo List',null,null,true,1
	    from ac_act_grps actgrp left join ac_res_grps acresgrp on true
	    where actgrp.name in ('ToDo Action Group') and acresgrp.name in ('ToDo List Access') on conflict do nothing;`)

	if toDoListPol.Error != nil {
		log.Fatal("Error")
	}

	toDoCreatePol := db.Exec(`INSERT INTO public.ac_policies(
		id, created_at, updated_at, deleted_at, ac_act_grp_id, ac_res_grp_id, owner_space_id, created_by, updated_by, name, description, path, opt_counter, is_predefined, type)
	    select '6YoEF0QogY4n4_4CBEHoyADibdHKRx7Z',now(),now(),null,actgrp.id,acresgrp.id,null,null,null,'ToDo Create','Ac Policy for ToDo Create',null,null,true,1
	    from ac_act_grps actgrp left join ac_res_grps acresgrp on true
	    where actgrp.name in ('ToDo Action Group') and acresgrp.name in ('ToDo Create Access') on conflict do nothing;`)

	if toDoCreatePol.Error != nil {
		log.Fatal("Error")
	}

	todoRemovePol := db.Exec(`INSERT INTO public.ac_policies(
		id, created_at, updated_at, deleted_at, ac_act_grp_id, ac_res_grp_id, owner_space_id, created_by, updated_by, name, description, path, opt_counter, is_predefined, type)
	    select 'VgdtYWT6lMD1UMcrTMRdWYbskU2rmJKo',now(),now(),null,actgrp.id,acresgrp.id,null,null,null,'ToDo Delete','Ac Policy for ToDo Delete',null,null,true,1
	    from ac_act_grps actgrp left join ac_res_grps acresgrp on true
	    where actgrp.name in ('ToDo Action Group') and acresgrp.name in ('ToDo Delete Access') on conflict do nothing;`)

	if todoRemovePol.Error != nil {
		log.Fatal("Error")
	}

	//seeding ac_pol_grps  table - secondary table for ac_policies
	todoListPolGrp := db.Exec(`INSERT INTO public.ac_pol_grps(
		id, created_at, updated_at, deleted_at, owner_space_id, description, name, opt_counter, is_predefined, type, entity_type, display_name, entity_types)
	    select '_w9SN6jrC2NNZx7xjNwRo2V5VoeRNHIk',now(),now(),null,null,'Ac Policy Group for ToDo List Access','ToDo List Access',null,true,1,null,'ToDo List Access','{1,2,3}' on conflict do nothing;`)

	if todoListPolGrp.Error != nil {
		log.Fatal("Error")
	}

	todoCreatePolGrp := db.Exec(`INSERT INTO public.ac_pol_grps(
		id, created_at, updated_at, deleted_at, owner_space_id, description, name, opt_counter, is_predefined, type, entity_type, display_name, entity_types)
	    select 'NIH-2a-xN-DdrR3UoynoofGoGb9jlOhQ',now(),now(),null,null,'Ac Policy Group for ToDo Create Access','ToDo Create Access',null,true,1,null,'ToDo Create Access','{1,2,3}' on conflict do nothing;`)

	if todoCreatePolGrp.Error != nil {
		log.Fatal("Error")
	}

	todoRemovePolGrp := db.Exec(`INSERT INTO public.ac_pol_grps(
		id, created_at, updated_at, deleted_at, owner_space_id, description, name, opt_counter, is_predefined, type, entity_type, display_name, entity_types)
	    select 'aIh9GZrLeIp8VtlWBXv9GZt-rJlOw9hS',now(),now(),null,null,'Ac Policy Group for ToDo Delete Access','ToDo Delete Access',null,true,1,null,'ToDo Delete Access','{1,2,3}' on conflict do nothing;`)

	if todoRemovePolGrp.Error != nil {
		log.Fatal("Error")
	}

	//seeding pol_gp_policies - ac_pol_grps and ac_policies bridge table
	listPolGrpPol := db.Exec(`INSERT INTO public.pol_gp_policies(
		id, created_at, updated_at, deleted_at, ac_policy_id, ac_pol_grp_id, opt_counter)
		select 'zJq6MM5lGlKiD9CurdO1rm7HJ1o25qrE',now(),now(),null,acpol.id,acpolgrp.id,null
		from ac_policies acpol left join ac_pol_grps acpolgrp on true
		where acpol.name in ('ToDo List') and acpolgrp.name in ('ToDo List Access') on conflict do nothing;`)

	if listPolGrpPol.Error != nil {
		log.Fatal("Error")
	}

	createPolGrpPol := db.Exec(`INSERT INTO public.pol_gp_policies(
		id, created_at, updated_at, deleted_at, ac_policy_id, ac_pol_grp_id, opt_counter)
		select 'iR7VMVc_gpaG2KnTP687qKrXWLLJKsrf',now(),now(),null,acpol.id,acpolgrp.id,null
		from ac_policies acpol left join ac_pol_grps acpolgrp on true
		where acpol.name in ('ToDo Create') and acpolgrp.name in ('ToDo Create Access') on conflict do nothing;`)

	if createPolGrpPol.Error != nil {
		log.Fatal("Error")
	}

	deletePolGrpPol := db.Exec(`INSERT INTO public.pol_gp_policies(
		id, created_at, updated_at, deleted_at, ac_policy_id, ac_pol_grp_id, opt_counter)
		select 'ktKzjpdFlmXm4Rxz45jP78VtluzhWQ2F',now(),now(),null,acpol.id,acpolgrp.id,null
		from ac_policies acpol left join ac_pol_grps acpolgrp on true
		where acpol.name in ('ToDo Delete') and acpolgrp.name in ('ToDo Delete Access') on conflict do nothing;`)

	if deletePolGrpPol.Error != nil {
		log.Fatal("Error")
	}

	//seeding permissions table - permissions basically API usage access
	listPer := db.Exec(`INSERT INTO public.ac_permissions(
		id, created_at, updated_at, deleted_at, description, name, opt_counter, is_predefined, type, display_name)
	    select '2MCxOl_RljjXhyc6hlWWhGdn5bmrCkpK',now(),now(),null,'Ac Permission for ToDo List And Create','ToDo List And Create Permission',null,true,1,'ToDo List And Create Permission'
	    on conflict do nothing;`)

	if listPer.Error != nil {
		log.Fatal("Error")
	}

	deletePer := db.Exec(`INSERT INTO public.ac_permissions(
		id, created_at, updated_at, deleted_at, description, name, opt_counter, is_predefined, type, display_name)
	    select 'SiQJGFqz1u6ZEQSKX8sR2wuJ1cxZNDX7',now(),now(),null,'Ac Permission for ToDo Delete Permission','ToDo Delete Permission',null,true,1,'ToDo Delete Permission'
	    on conflict do nothing;`)

	if deletePer.Error != nil {
		log.Fatal("Error")
	}

	listPerPolGrp := db.Exec(`INSERT INTO public.per_pol_grps(
	    id, created_at, updated_at, deleted_at, ac_permission_id, ac_pol_grp_id)
	    select 'Wihu1Ff2GxC7C_B6WEYKbVdqZaVjCwK',now(),now(),null,acper.id,acpolgrp.id from ac_permissions acper
	    left join ac_pol_grps acpolgrp on true
	    where acper.name in ('ToDo List And Create Permission') and acpolgrp.name in ('ToDo Create Access','ToDo List Access')
		on conflict do nothing;`)

	if listPerPolGrp.Error != nil {
		log.Fatal("Error")
	}

	createPerPolGrp := db.Exec(`INSERT INTO public.per_pol_grps(
	    id, created_at, updated_at, deleted_at, ac_permission_id, ac_pol_grp_id)
	    select 'H45k5fyzdj3I2Kab4_Rn64na_0-WqvUF',now(),now(),null,acper.id,acpolgrp.id from ac_permissions acper
	    left join ac_pol_grps acpolgrp on true
	    where acper.name in ('ToDo List And Create Permission') and acpolgrp.name in ('ToDo Create Access')
		on conflict do nothing;`)

	if createPerPolGrp.Error != nil {
		log.Fatal("Error")
	}

	//seeding per_pol_grps - ac_permissions and  policies connection table
	deletePerPolGrp := db.Exec(`INSERT INTO public.per_pol_grps(
	    id, created_at, updated_at, deleted_at, ac_permission_id, ac_pol_grp_id)
	    select 'U445o3QnvrzBE2PohtajqDp3hA_PSvcg',now(),now(),null,acper.id,acpolgrp.id from ac_permissions acper
	    left join ac_pol_grps acpolgrp on true
	    where acper.name in ('ToDo Delete Permission') and acpolgrp.name in ('ToDo Delete Access')
		on conflict do nothing;`)

	if deletePerPolGrp.Error != nil {
		log.Fatal("Error")
	}
}
