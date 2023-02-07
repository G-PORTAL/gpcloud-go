# GPCloud Golang Client

This is the official GPCloud Golang client. Please raise an issue if you have found any problems or having questions.

### Recommendations

- Golang 1.18 or higher

### Example usage

```go
package main

import (
	"context"
	"log"

	authv1 "buf.build/gen/go/gportal/gportal-cloud/protocolbuffers/go/gpcloud/api/auth/v1"
	cloudv1 "buf.build/gen/go/gportal/gportal-cloud/protocolbuffers/go/gpcloud/api/cloud/v1"
	"github.com/G-PORTAL/gpcloud-go/pkg/gpcloud/client"
	"github.com/G-PORTAL/gpcloud-go/pkg/gpcloud/client/auth"
)

func main() {
	conn, err := client.NewClient(
		// For getting your own client ID and client Secret please ask support
		&auth.ProviderKeycloakUserPassword{
			ClientID:     "my-custom-client-id",
			ClientSecret: "my-custom-client-secret",
			Username:     "example@gpcloud.customer",
			Password:     "password123",
		},
	)
	if err != nil {
		log.Fatal("failed to create client:\n", err)
	}

	ctx := context.Background()
	user, err := conn.AuthClient().GetUser(ctx, &authv1.GetUserRequest{})
	if err != nil {
		log.Fatal("failed to fetch user information:\n", err)
	}
	log.Println("User ID: ", user.User.Id)

	projects, err := conn.CloudClient().ListProjects(ctx, &cloudv1.ListProjectsRequest{})
	if err != nil {
		log.Fatal("failed to fetch project list: \n", err)
	}
	for _, project := range projects.Projects {
		log.Println("Project ID: ", project.Id)
	}
}

```