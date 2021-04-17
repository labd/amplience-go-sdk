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

	Body                 map[string]interface{}
	Version              int             `json:"version"`
	Label                string          `json:"label"`
	Status               string          `json:"status"`
	CreatedBy            string          `json:"createdBy"`
	Locale               string          `json:"locale"`
	CreatedDate          *time.Time      `json:"createdDate"`
	LastModifiedBy       string          `json:"lastModifiedBy"`
	LastModifiedDate     *time.Time      `json:"lastModifiedDate"`
	LastPublishedVersion int             `json:"lastPublishedVersion"`
	LastPublishedDate    *time.Time      `json:"lastPublishedDate"`
	DeliveryID           string          `json:"deliveryId"`
	Links                map[string]Link `json:"_links"`
}

type ContentItemResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []ContentItem
}

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

func (client *Client) ContentItemGet() {

}

func (client *Client) ContentItemUpdate() {

}

func (client *Client) ContentItemCreate() {

}

func (client *Client) ContentItemList(repositoryId string) (ContentItemResults, error) {
	result := ContentItemResults{}
	endpoint := fmt.Sprintf("/content-repositories/%s/content-items", repositoryId)

	err := client.request(http.MethodGet, endpoint, &result)
	return result, err

}

func (client *Client) ContentItemArchive() {

}

func (client *Client) ContentItemUnarchive() {

}

func (client *Client) ContentItemListHistory() {

}
