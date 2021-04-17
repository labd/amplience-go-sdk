# Amplience GO SDK

GO SDK for [Amplience](https://amplience.com/).


## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/labd/amplience-go-sdk/content"
)

func main() {

	client, err := content.NewClient(&content.ClientConfig{
		ClientID:     "<my-client-id>",
		ClientSecret: "<my-client-secret>",
	})

	results, err := client.HubList()
	for _, hub := range results.Items {
		log.Println(hub.ID)
	}
	hub, err := client.HubGet("<my-hub-id>")

  results, err = client.ContentRepositoryList("<my-hub-id>")
	for _, repository := range results.Items {
		log.Println(repository.ID)
	}

  repository, err := client.ContentRepositoryGet("<my-repository-id>")
	results, err = client.ContentItemList(repository.ID)
	for _, item := range results.Items {
		log.Println(item.ID)
	}

}
```
