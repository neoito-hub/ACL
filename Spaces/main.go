package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/neoito-hub/ACL-Block/spaces/common_services"
	gen "github.com/neoito-hub/ACL-Block/spaces/gen/go/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Service struct {
	gen.UnimplementedSpacesProxyServer
}

func (g *Service) SpacesCallService(ctx context.Context, request *gen.SpacesRequest) (*gen.SpacesReply, error) {

	// parsing body.
	var userIdValues []string
	var userID string
	var urlValues []string
	var url string
	var userNameValues []string
	var userName string
	var spaceIdValues []string
	var isOwner string
	var isOwnerValues []string
	var spaceId string

	//decoding url and userID from request context
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userIdValues = md.Get("user-id")
		urlValues = md.Get("url")
		userNameValues = md.Get("user-name")
		spaceIdValues = md.Get("space-id")
		isOwnerValues = md.Get("is-owner")

	}

	if len(userIdValues) > 0 {
		userID = userIdValues[0]
	}
	if len(urlValues) > 0 {
		url = urlValues[0]
	}
	if len(userNameValues) > 0 {
		userName = urlValues[0]
	}

	if len(spaceIdValues) > 0 {
		spaceId = spaceIdValues[0]
	}

	if len(isOwnerValues) > 0 {
		isOwner = isOwnerValues[0]
	}

	invReply := InvokeSpacesFunction(funcs, common_services.HandlerPayload{Url: url, RequestBody: request.Body, UserID: userID, UserName: userName, Queryparams: request.Queryparams, SpaceID: spaceId, IsOwner: isOwner})

	// result := timed(InvokeFunctionPayload{UserID: userID, Url: url, RequestBody: request.Body})

	fmt.Printf("invite reply is %v", invReply)

	// // converting pb struct to Map Interface.
	// body := request.Body.AsMap()

	// // convert map to json string
	// jsonString, _ := json.Marshal(body)
	// s := (string(jsonString))
	// log.Println(s)
	// log.Println(reflect.TypeOf(s))

	// // converting json string to struct
	// as_struct := Body{}
	// err := json.Unmarshal([]byte(s), &as_struct)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// log.Println(as_struct)
	// log.Println("#####################################")

	if err := request.Validate(); err != nil {
		return nil, err
	}

	// strWrapper := wrapperspb.String("foo")
	// strAny, _ := anypb.New(strWrapper)

	return &gen.SpacesReply{
		Err:    invReply.Err,
		Data:   invReply.Data,
		Status: int32(invReply.Status),
	}, nil
}

func main() {
	// 	// Load env vars
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading .env file: %v", envErr)
	}

	//Load func map
	funcloadErr := loadFuncs()
	if funcloadErr != nil {
		log.Fatalf("Error loading function map %v", funcloadErr)
	}

	//Initialise common db object for grpc handlers invocation
	DBInit()
	// defer CloseDbCOnn()

	// UpdateAcResources()

	// create new gRPC server
	server := grpc.NewServer()
	// register the GreeterServerImpl on the gRPC server
	gen.RegisterSpacesProxyServer(server, &Service{})
	// start listening on port :8080 for a tcp connection

	fmt.Println("######################################")
	log.Println(fmt.Sprintf("Starting SPACES App %s", os.Getenv("SPACES_PORT")))
	fmt.Println("######################################")

	if l, err := net.Listen("tcp", os.Getenv("SPACES_PORT")); err != nil {
		log.Fatal("error in listening on port :5000", err)
	} else {
		// the start gRPC server
		if err := server.Serve(l); err != nil {
			log.Fatal("unable to start server", err)
		}
	}

}
