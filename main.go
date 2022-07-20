package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gravitational/teleport/api/client"
	"github.com/gravitational/teleport/api/types"
)

var (
	proxyParameter             *string
	userIdParameter            *string
	rolesParameter             *string
	resourceTypesParameter     *string
	resourceLabelsParameter    *string
	identityFilePathParameter  *string
	waitForUserParameter       *bool
	timeToWaitSecondsParameter *time.Duration
	approveSubmittedRequest    *bool
)

func init() {
	proxyParameter = flag.String("proxy", "", "Teleport Cluster Proxy Address (teleport.example.com:443)")
	userIdParameter = flag.String("user", "", "Teleport user")
	rolesParameter = flag.String("roles", "", "Teleport comma delimited roles to submit request for")
	identityFilePathParameter = flag.String("identity", "", "Identity file for submitting and approving request")
	waitForUserParameter = flag.Bool("waitforuser", true, "Wait for user to have an entry")
	timeToWaitSecondsParameter = flag.Duration("time", 5, "Seconds to wait in-between waiting for user to appear")
	approveSubmittedRequest = flag.Bool("approverequest", true, "Submit a approval review for the requested roles")
}

func main() {

	flag.Parse()

	if len(*proxyParameter) == 0 || len(*userIdParameter) == 0 || len(*rolesParameter) == 0 {
		log.Printf("Proxy, user id and roles required.")
		flag.PrintDefaults()

		return
	}

	roles := strings.Split(*rolesParameter, ",")

	log.Printf("Submitting access request for %s on %s for roles %s", *proxyParameter, *userIdParameter, roles)

	ctx := context.Background()

	clt, err := client.New(ctx, client.Config{
		Addrs: []string{

			*proxyParameter,
		},
		Credentials: getCredentials(*identityFilePathParameter),
	})

	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}

	defer clt.Close()
	resp, err := clt.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping server: %v", err)
		return
	} else {
		log.Printf("Connected to server: %s", resp)
	}

	userNotFound := true
	for userNotFound == true {
		users, usersErr := clt.GetUsers(false)
		if usersErr != nil {
			log.Fatalf("users retrieval error: %v", err)
			return
		}
		for i := 0; i < len(users); i++ {
			if users[i].GetName() == *userIdParameter {
				userNotFound = false
				break
			}
		}
		if userNotFound == true {

			if *waitForUserParameter == true {

				log.Printf("User %s not found, continuing wait for user to become available ", *userIdParameter)
				var timeseconds = *timeToWaitSecondsParameter * time.Second
				time.Sleep(timeseconds)
			} else {
				log.Fatalf("User %s not found, exitting.", *userIdParameter)
				return
			}
		}
	}

	requestId := uuid.New().String()
	req, err := types.NewAccessRequestWithResources(requestId, *userIdParameter, roles, nil)

	err2 := clt.CreateAccessRequest(ctx, req)
	if err2 != nil {
		log.Fatalf("failed to create access request: %v", err2)
		return
	} else {

		if *approveSubmittedRequest == true {
			accessRequest, errAccessRequest := clt.SubmitAccessReview(ctx, types.AccessReviewSubmission{
				RequestID: requestId,
				Review: types.AccessReview{
					ProposedState: types.RequestState_APPROVED,
					Reason:        "Approved",
					Created:       time.Now(),
				},
			})
			log.Printf("Access request state: %v", accessRequest.GetState())
			if errAccessRequest != nil {
				log.Fatalf("failed to update access request: %v", errAccessRequest)
				return
			}
		}
	}

}

func getCredentials(identityPath string) []client.Credentials {
	if len(identityPath) != 0 {
		return []client.Credentials{client.LoadIdentityFile(identityPath)}
	} else {
		return []client.Credentials{client.LoadProfile("", "")}
	}
}
