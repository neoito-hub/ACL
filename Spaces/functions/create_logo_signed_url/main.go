package create_logo_signed_url

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aidarkhanov/nanoid"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/neoito-hub/ACL-Block/spaces/common_services"
)

func Handler(payload common_services.HandlerPayload) common_services.HandlerResponse {
	// Validating request body and method
	// b, validateErr := ValidateRequest(w, r)
	// if validateErr != nil {
	// 	fmt.Printf("Error: %v\n", validateErr)
	// 	RespondWithError(w, http.StatusBadRequest, validateErr.Error())

	// 	return
	// }

	// // Validating and retreving user id from user access token
	// _, shieldVerifyError := VerifyAndGetUser(w, r)
	// if shieldVerifyError != nil {
	// 	fmt.Printf("shieldVerifyError: %v\n", shieldVerifyError)
	// 	RespondWithError(w, http.StatusUnauthorized, shieldVerifyError.Error())

	// 	return
	// }

	var b RequestObject
	var handlerResp common_services.HandlerResponse

	if err := json.Unmarshal([]byte(payload.RequestBody), &b); err != nil {
		handlerResp = common_services.BuildErrorResponse(true, "Invalid Request Payload", Response{}, http.StatusBadRequest)
		return handlerResp

	}

	// Load the bucket env
	s3BlockBucket := os.Getenv("S3_BLOX_LOGO_BUCKET")
	s3BlockBucketRegion := os.Getenv("S3_BLOX_LOGO_BUCKET_REGION")

	// Prepare the S3 request so a signature can be generated
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3BlockBucketRegion)},
	)
	if err != nil {
		fmt.Println("Failed to create a new session: ", err)
		// RespondWithError(w, http.StatusInternalServerError, "Error creating session")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating session", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	svc := s3.New(sess)

	Key := nanoid.New() + time.Now().GoString() + "." + b.FileExtension
	awsKeyString := aws.String(Key)

	resp, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		// ACL:    aws.String("public-read"),
		Bucket: aws.String(s3BlockBucket),
		Key:    awsKeyString,
	})

	var expiryTimeValue time.Duration = 60 // in minutes

	// Create the pre-signed url with an expiry
	url, err := resp.Presign(expiryTimeValue * time.Minute)

	if err != nil {
		fmt.Println("Failed to generate a pre-signed url: ", err)
		// RespondWithError(w, http.StatusInternalServerError, "Error creating pre-signed url")

		// return

		handlerResp = common_services.BuildErrorResponse(true, "Error creating pre-signed url", Response{}, http.StatusInternalServerError)
		return handlerResp
	}

	// RespondWithJSON(w, http.StatusOK, ReturnObject{URL: url, Key: Key})

	handlerResp = common_services.BuildResponse(false, "pre-signed url generated successfully", ReturnObject{URL: url, Key: Key}, http.StatusOK)
	return handlerResp
}
