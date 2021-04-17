package content

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ContentTypeReference struct {
	HubContentTypeID string `json:"hubContentTypeId"`
	ContentTypeUri   string `json:"contentTypeUri"`
}

type ContentRepository struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Label        string                 `json:"label"`
	Status       string                 `json:"status"`
	Type         string                 `json:"type"`
	ContentTypes []ContentTypeReference `json:contentTypes"`
}

type ContentRepositoryResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []ContentRepository
}

func (r *ContentRepositoryResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["content-repositories"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (client *Client) ContentRepositoryGet(id string) (ContentRepository, error) {
	result := ContentRepository{}
	endpoint := fmt.Sprintf("/content-repositories/%s", id)
	err := client.request(http.MethodGet, endpoint, &result)
	return result, err
}

func (client *Client) ContentRepositoryCreate() {

}

func (client *Client) ContentRepositoryList(hubId string) (ContentRepositoryResults, error) {
	result := ContentRepositoryResults{}
	endpoint := fmt.Sprintf("/hubs/%s/content-repositories", hubId)
	err := client.request(http.MethodGet, endpoint, &result)
	return result, err
}

func (client *Client) ContentRepositoryUpdate() {

}

func (client *Client) ContentRepositoryFind() {

}

func (client *Client) ContentRepositoryShare() {

}

func (client *Client) ContentRepositoryAssignContentType(repositoryId string) {

}

func (client *Client) ContentRepositoryAssignRemoveContentType(repositoryId string, typeId string) {

}

func (client *Client) ContentRepositoryAssignFeature(repositoryId string) {

}

func (client *Client) ContentRepositoryRemoveFeature(repositoryId string) {

}
