module github.com/neoito-hub/ACL-Block/spaces

go 1.18

require (
	github.com/aidarkhanov/nanoid v1.0.8
	//github.com/appblocks-hub/appblocks-datamodels-backend/models v0.0.0-20230724044841-f9c9dcaa9015
	github.com/aws/aws-sdk-go v1.49.14
	github.com/envoyproxy/protoc-gen-validate v1.0.2
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	google.golang.org/genproto v0.0.0-20240102182953-50ed04b92917
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
	gorm.io/datatypes v1.2.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

require github.com/jackc/pgx/v5 v5.5.1 // indirect

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgtype v1.14.0
	github.com/jackc/pgx/v4 v4.18.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	gorm.io/driver/mysql v1.5.2 // indirect
)

require github.com/neoito-hub/ACL-Block/Data-Models v0.0.0-00010101000000-000000000000

require (
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/lyft/protoc-gen-star/v2 v2.0.3 // indirect
	github.com/spf13/afero v1.3.3 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
)

replace github.com/neoito-hub/ACL-Block/Data-Models => ./Data-Models

