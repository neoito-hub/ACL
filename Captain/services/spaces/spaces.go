package spaces

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/neoito-hub/ACL-Block/captain/common_services"
	pb "github.com/neoito-hub/ACL-Block/captain/gen/go/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// type RequestObject struct {
// 	data map[string]
// }

type RedirectResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

func InvokeGRPC(w http.ResponseWriter, r *http.Request, shieldUser common_services.ContextData, routeData RouteData) {
	urlMap := make(map[string]string)
	r.ParseForm()

	for key, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			urlMap[key] = value
		}
	}

	//convert request body to string for passing to server
	body, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}
	stringBody := string(body)

	conn, err := grpc.Dial(os.Getenv("SPACES_DIAL_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSpacesProxyClient(conn)

	header := metadata.New(map[string]string{
		"user-id": shieldUser.UserID, "url": routeData.Url, "host": routeData.Host, "user-name": shieldUser.UserName, "space-id": shieldUser.SpaceID, "is-owner": strconv.FormatBool(shieldUser.IsOwner),
	})

	ctx := metadata.NewOutgoingContext(context.Background(), header)

	reply, err := client.SpacesCallService(ctx, &pb.SpacesRequest{Body: stringBody, Queryparams: urlMap})
	if err != nil {
		fmt.Printf("error is", err)
		respondWithJSON(w, 500, "")
	}
	log.Printf("data: %v", reply)

	if reply.Status == http.StatusTemporaryRedirect {
		respondWithRedirect(w, int(reply.Status), reply.Data, r)
	}

	respondWithJSON(w, int(reply.Status), reply.Data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}

func respondWithRedirect(w http.ResponseWriter, code int, payload string, r *http.Request) {
	type RedirectResponse struct {
		RedirectUrl string `json:"rediDatarect_url"`
	}

	var payloadObj RedirectResponse

	if err := json.Unmarshal([]byte(payload), &payloadObj); err != nil {
		respondWithJSON(w, 500, "")
		return
	}

	fmt.Printf("redirect url is %v", payloadObj.RedirectUrl)

	r.Header.Del("Authorization")

	http.Redirect(w, r, payloadObj.RedirectUrl, code)
}
