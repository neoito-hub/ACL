package roles_list_available_permissions

import (
	"fmt"
	"os"

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

	dbinf.Host = os.Getenv("SPACES_POSTGRES_HOST")
	dbinf.User = os.Getenv("SPACES_POSTGRES_USER")
	dbinf.Password = os.Getenv("SPACES_POSTGRES_PASSWORD")
	dbinf.Name = os.Getenv("SPACES_POSTGRES_NAME")
	dbinf.Port = os.Getenv("SPACES_POSTGRES_PORT")
	dbinf.Sslmode = os.Getenv("SPACES_POSTGRES_SSLMODE")
	dbinf.Timezone = os.Getenv("SPACES_POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbinf.Host, dbinf.User, dbinf.Password, dbinf.Name, dbinf.Port, dbinf.Sslmode, dbinf.Timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("DB connection err")
	}

	return db
}
