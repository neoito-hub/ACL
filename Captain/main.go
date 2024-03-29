package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neoito-hub/ACL-Block/captain/common_services"
	"github.com/neoito-hub/ACL-Block/captain/gateway"
	"github.com/neoito-hub/ACL-Block/captain/services/spaces"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var resourcesMap = make(map[string]common_services.Resources)

func LoadResources() {
	type Data struct {
		EntiyName       string `json:"entity_name"`
		IsAuthorised    int    `json:"is_authorised"`
		IsAuthenticated int    `json:"is_authenticated"`
		Path            string `json:"path"`
		FunctionName    string `json:"function_name"`
	}

	var resourcesData []Data

	res := db.Raw("select entity_name,is_authorised,is_authenticated,path,function_name from ac_resources ac ").Scan(&resourcesData)

	if res.Error != nil {
		log.Fatal("Could not load resources from db")
	}

	for _, v := range resourcesData {
		resourcesMap[v.FunctionName] = common_services.Resources{EntiyName: v.EntiyName, IsAuthorised: v.IsAuthorised, IsAuthenticated: v.IsAuthenticated}

	}

}

func DBInit() {
	dbinf := &common_services.DBInfo{}
	var dbErr error

	dbinf.Host = os.Getenv("DB_POSTGRES_HOST")
	dbinf.User = os.Getenv("DB_POSTGRES_USER")
	dbinf.Password = os.Getenv("DB_POSTGRES_PASSWORD")
	dbinf.Name = os.Getenv("DB_POSTGRES_NAME")
	dbinf.Port = os.Getenv("DB_POSTGRES_PORT")
	dbinf.Sslmode = os.Getenv("DB_POSTGRES_SSLMODE")
	dbinf.Timezone = os.Getenv("DB_POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbinf.Host, dbinf.User, dbinf.Password, dbinf.Name, dbinf.Port, dbinf.Sslmode, dbinf.Timezone)
	db, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if dbErr != nil {
		panic("DB connection err")
	}
}

func CloseDbCOnn() {
	//closing connection to db
	sqlDB, dberr := db.DB()
	if dberr != nil {
		log.Fatalln(dberr)
	}

	defer sqlDB.Close()

}

func SpacesHandler(c *gin.Context) {

	w := c.Writer
	r := c.Request

	urlFragments := strings.Split(r.URL.Path, "/")

	actionName := urlFragments[len(urlFragments)-1]
	functionName := urlFragments[len(urlFragments)-2]

	spaceID := r.Header.Get("space_id")

	mdlwrErr, shieldUser := gateway.Call(w, r, db, resourcesMap, actionName, functionName, spaceID)
	if mdlwrErr != nil {
		return
	}
	spaces.InvokeGRPC(w, r, shieldUser, spaces.RouteData{Url: r.URL.Path, Host: r.Host})

}

func AuthHandler(c *gin.Context) {

	w := c.Writer
	r := c.Request

	actionName := r.URL.Query().Get("action-name")
	resourceName := r.URL.Query().Get("resource-name")
	spaceID := r.URL.Query().Get("space-id")

	mdlwrErr, shieldUser := gateway.Call(w, r, db, resourcesMap, actionName, resourceName, spaceID)
	if mdlwrErr != nil {
		return
	}

	common_services.RespondWithJSON(w, http.StatusOK, shieldUser)

}

func main() {

	// Load env vars
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading .env file: %v", envErr)
	}

	//Initialise common db object for grpc handlers invocation
	DBInit()
	// defer CloseDbCOnn()

	LoadResources()

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}

	r := gin.Default()

	r.POST("/api/spaces/*action", SpacesHandler)
	r.GET("/api/spaces/*action", SpacesHandler)
	r.PUT("/api/spaces/*action", SpacesHandler)
	r.DELETE("/api/spaces/*action", SpacesHandler)

	r.GET("/api/auth/getUser", AuthHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowedHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization", "space_id"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}).Handler(r)

	fmt.Println("\n######################################")
	log.Println(fmt.Sprintf("Starting App%s", os.Getenv("GATEWAY_PORT")))
	fmt.Println("######################################")
	log.Fatal(http.ListenAndServe(os.Getenv("GATEWAY_PORT"), handler))

}
