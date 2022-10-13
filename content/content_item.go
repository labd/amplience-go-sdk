package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ContentItem struct {
	ID                  string `json:"id"`
	ContentRepositoryID string `json:"contentRepositoryId"`
	FolderID            string `json:"folderId"`

	Body                 map[string]interface{} `json:"body"`
	Version              int                    `json:"version"`
	Label                string                 `json:"label"`
	Status               string                 `json:"status"`
	CreatedBy            string                 `json:"createdBy"`
	Locale               string                 `json:"locale"`
	CreatedDate          *time.Time             `json:"createdDate"`
	LastModifiedBy       string                 `json:"lastModifiedBy"`
	LastModifiedDate     *time.Time             `json:"lastModifiedDate"`
	LastPublishedVersion int                    `json:"lastPublishedVersion"`
	LastPublishedDate    *time.Time             `json:"lastPublishedDate"`
	DeliveryID           string                 `json:"deliveryId"`
	Links                map[string]Link        `json:"_links"`
}

type ContentItemVersionHistory struct {
	HistoryEventID string            `json:"historyEventId"`
	ContentItemId  string            `json:"contentItemId"`
	Version        int               `json:"version"`
	CreatedDate    *time.Time        `json:"createdDate"`
	CreatedBy      string            `json:"createdBy"`
	Action         ContentItemAction `json:"action"`
}

type ContentItemInput struct {
	Body     map[string]interface{} `json:"body"`
	Label    string                 `json:"label"`
	FolderID string                 `json:"folderId"`
	Locale   string                 `json:"locale"`
}

type ContentItemAction struct {
	Code string                 `json:"code"`
	Data map[string]interface{} `json:"data"`
}

// ContentItemResults is returned by the ContentItemList func
type ContentItemResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []ContentItem
}

// UnmarshalJSON is a custom unmarshaller for the embedded content
func (r *ContentItemResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["content-items"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

// ContentItemResults is returned by the ContentItemList func
type ContentItemVersionHistoryResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []ContentItemVersionHistory
}

// UnmarshalJSON is a custom unmarshaller for the embedded content
func (r *ContentItemVersionHistoryResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["content-item-version-history"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

// ContentItemCreate creates a new content item
func (client *Client) ContentItemCreate(repositoryID string, input ContentItemInput) (ContentItem, error) {
	result := ContentItem{}
	body, err := json.Marshal(input)
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/content-repositories/%s/content-items", repositoryID)
	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

// ContentItemGet returns the content item with the given id
func (client *Client) ContentItemGet(id string) (ContentItem, error) {
	endpoint := fmt.Sprintf("/content-items/%s", id)
	result := ContentItem{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

// ContentItemUpdate updates a Content Item. Please note that deliveryKey can
// only be set when Content Delivery 2 is enabled.
func (client *Client) ContentItemUpdate(current ContentItem, input ContentItemInput) (ContentItem, error) {
	result := ContentItem{}

	body, err := createUpdatePatch(
		ContentItemInput{
			Body:     current.Body,
			Label:    current.Label,
			FolderID: current.FolderID,
			Locale:   current.Locale,
		},
		input)

	if body == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/content-items/%s", current.ID)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

// ContentItemList lists all of the Content Items within the given Content
// Repository
func (client *Client) ContentItemList(repositoryID string, parameters ContentItemPaginationParameters) (ContentItemResults, error) {
	result := ContentItemResults{}
	endpoint := fmt.Sprintf("/content-repositories/%s/content-items?%s", repositoryID, ContentItemPaginationQueryString(parameters))

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentItemGetAll(hubID string, status ContentStatus) ([]ContentItem, error) {
	parameters := ContentItemPaginationParameters{}
	parameters.Status = status

	response, err := client.ContentItemList(hubID, parameters)

	var result []ContentItem
	result = append(result, response.Items...)

	for parameters.Page < response.Page.TotalPages-1 {
		parameters.Page++
		response, err := client.ContentItemList(hubID, parameters)
		if err != nil {
			break
		}
		result = append(result, response.Items...)
	}

	return result, err
}

// ContentItemArchive archives a content item
func (client *Client) ContentItemArchive(id string, version int) (ContentItem, error) {
	result := ContentItem{}
	endpoint := fmt.Sprintf("/content-types/%s/archive", id)

	body, err := json.Marshal(ArchiveInput{Version: version})
	if err != nil {
		return result, err
	}

	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

// ContentItemUnarchive unarchives a content item
func (client *Client) ContentItemUnarchive(id string, version int) (ContentItem, error) {
	result := ContentItem{}
	endpoint := fmt.Sprintf("/content-types/%s/unarchive", id)

	body, err := json.Marshal(ArchiveInput{Version: version})
	if err != nil {
		return result, err
	}

	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

// ContentItemListHistory list history of this item
func (client *Client) ContentItemListHistory(id string, version int) (ContentItemVersionHistoryResults, error) {
	endpoint := fmt.Sprintf("/content-items/%s/versions/%d/history", id, version)
	result := ContentItemVersionHistoryResults{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err

}
