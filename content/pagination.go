package content

import (
	"fmt"
	"net/url"
	"strings"
)

type PaginationParameters struct {
	Page int
	Size int
	Sort string
}

func PaginationQueryString(parameters PaginationParameters) string {

	q := url.Values{}
	q.Add("page", fmt.Sprintf("%v", parameters.Page))
	if parameters.Size > 0 {
		q.Add("size", fmt.Sprintf("%v", parameters.Size))
	}
	if parameters.Sort != "" {
		q.Add("sort", parameters.Sort)
	}

	return q.Encode()
}

type ContentStatus string

const (
	StatusAny      ContentStatus = ""
	StatusActive   ContentStatus = "ACTIVE"
	StatusArchived ContentStatus = "ARCHIVED"
)

type StatusPaginationParameters struct {
	Page   int
	Size   int
	Sort   string
	Status ContentStatus
}

func StatusPaginationQueryString(parameters StatusPaginationParameters) string {

	q := url.Values{}
	q.Add("page", fmt.Sprintf("%v", parameters.Page))
	if parameters.Size > 0 {
		q.Add("size", fmt.Sprintf("%v", parameters.Size))
	}
	if parameters.Status != StatusAny {
		q.Add("status", string(parameters.Status))
	}
	if parameters.Sort != "" {
		q.Add("sort", parameters.Sort)
	}

	return q.Encode()
}

type Projection string

const (
	ProjectionAny   Projection = ""
	ProjectionBasic Projection = "basic"
)

type ContentItemPaginationParameters struct {
	Page                        int
	Size                        int
	Sort                        string
	Status                      ContentStatus
	FolderId                    string
	Projection                  Projection
	ExcludeHierarchicalChildren bool
}

func ContentItemPaginationQueryString(parameters ContentItemPaginationParameters) string {

	q := url.Values{}
	q.Add("page", fmt.Sprintf("%v", parameters.Page))

	if parameters.Size > 0 {
		q.Add("size", fmt.Sprintf("%v", parameters.Size))
	} else if parameters.Projection != ProjectionBasic {
		// The default is 20, but if you're json response is too big, the Amplience API will return an error.
		// Therefore, to be safe we set a default of 6.
		q.Add("size", "6")
	}
	if parameters.Status != StatusAny {
		q.Add("status", string(parameters.Status))
	}
	if parameters.Sort != "" {
		q.Add("sort", parameters.Sort)
	}
	if parameters.ExcludeHierarchicalChildren {
		q.Add("excludeHierarchicalChildren", "true")
	}
	if parameters.FolderId != "" {
		q.Add("folderId", parameters.FolderId)
	}

	return strings.Replace(q.Encode(), "%2C", ",", 1)
}
