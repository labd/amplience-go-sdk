# Amplience Go SDK

A Go client for the [Amplience](https://amplience.com/) Dynamic Content
[Content Management API](https://amplience.com/developers/docs/apis/content-management-reference/).
It wraps the REST API in typed structs and methods, so Go programs can manage
hubs, content repositories, content types, content type schemas, content items,
folders, extensions, webhooks and Algolia search indexes without hand-building
HTTP requests. It handles OAuth2 client-credentials authentication and pagination
for you. Its main consumer is
[terraform-provider-amplience](https://github.com/labd/terraform-provider-amplience).

## Installation

```sh
go get github.com/labd/amplience-go-sdk
```

All functionality lives in the `content` package:

```go
import "github.com/labd/amplience-go-sdk/content"
```

## Requirements

- Go 1.27 or newer (see `go.mod`)
- An Amplience API client ID and secret with access to the Content Management API

## Usage

```go
package main

import (
	"log"
	"os"

	"github.com/labd/amplience-go-sdk/content"
)

func main() {
	client, err := content.NewClient(&content.ClientConfig{
		ClientID:     os.Getenv("AMPLIENCE_CLIENT_ID"),
		ClientSecret: os.Getenv("AMPLIENCE_CLIENT_SECRET"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// List the hubs this client can access (first page)
	hubs, err := client.HubList(content.PaginationParameters{})
	if err != nil {
		log.Fatal(err)
	}
	for _, hub := range hubs.Items {
		log.Println(hub.ID, hub.Name)
	}

	// Fetch all content repositories of a hub; *GetAll helpers follow pagination
	repositories, err := client.ContentRepositoryGetAll("<hub-id>")
	if err != nil {
		log.Fatal(err)
	}

	// List archived content items in the first repository
	items, err := client.ContentItemList(repositories[0].ID, content.ContentItemPaginationParameters{
		Status: content.StatusArchived,
		Size:   20,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range items.Items {
		log.Println(item.ID, item.Label, item.Version)
	}
}
```

Most resources follow the same pattern: `<Resource>Create`, `<Resource>Get`,
`<Resource>Update(current, input)`, `<Resource>List(parentID, params)` and
`<Resource>GetAll`, plus `Archive`/`Unarchive` or `Delete` where the API
supports them. See the [`content` package](./content) or
[pkg.go.dev](https://pkg.go.dev/github.com/labd/amplience-go-sdk/content)
for the full list.

### Errors

API errors are returned as `*content.ErrorResponse`, which carries the HTTP
status code and the error objects from the API:

```go
var apiErr *content.ErrorResponse
if errors.As(err, &apiErr) {
	log.Println(apiErr.StatusCode, apiErr.Errors)
}
```

## Configuration

`content.ClientConfig` accepts:

| Field          | Default                                 | Description                                    |
| -------------- | --------------------------------------- | ---------------------------------------------- |
| `ClientID`     |                                         | Amplience API client ID                        |
| `ClientSecret` |                                         | Amplience API client secret                    |
| `URL`          | `https://api.amplience.net/v2/content`  | Content Management API base URL                |
| `AuthURL`      | `https://auth.adis.ws/oauth/token`      | OAuth2 token endpoint                          |
| `HTTPClient`   | `http.DefaultClient`                    | Custom `*http.Client` used for all requests    |

The SDK does not read credentials from the environment itself. The only
environment variable it reads is `AMPLIENCE_DEBUG`: set it to any non-empty
value to log every request and response.

## Development

The repository uses [Task](https://taskfile.dev/):

```sh
task test      # go test -v ./...
task format    # go fmt ./...
task coverage  # tests with race detector and coverage report
```

CI (`.github/workflows/tests.yaml`) runs `go test -race ./...`, golangci-lint
and govulncheck.

To try changes against a real Amplience account, write a small `main.go` like
the usage example above and run it with `go run main.go`.

## Releasing

Changes are tracked with [changie](https://changie.dev/). Add a change file
with every user-facing change:

```sh
changie new
```

On every push to `main`,
[changie-release-action](https://github.com/labd/changie-release-action)
opens or updates a release PR that batches the unreleased changes into
`CHANGELOG.md`. Merging that PR tags the new version, which makes it available
through `go get`.

## License

MIT, see [LICENSE](./LICENSE).
