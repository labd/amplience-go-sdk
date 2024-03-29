package content

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ContentTypeReference struct {
	HubContentTypeID string `json:"hubContentTypeId"`
	ContentTypeURI   string `json:"contentTypeUri"`
}

type ContentRepository struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	Status       string `json:"status"`
	Type         string `json:"type"`
	HubID        string
	Links        map[string]Link        `json:"_links"`
	ContentTypes []ContentTypeReference `json:"contentTypes"`
}

func (r *ContentRepository) GetHub(client *Client) (Hub, error) {
	result := Hub{}
	err := client.request(http.MethodGet, r.Links["hub"].Href, nil, &result)
	return result, err
}

type ContentRepositoryInput struct {
	Name  string `json:"name"`
	Label string `json:"label"`
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

// ContentRepositoryGet returns a ContentRepository for the given id
func (client *Client) ContentRepositoryGet(id string) (ContentRepository, error) {
	result := ContentRepository{}
	endpoint := fmt.Sprintf("/content-repositories/%s", id)
	err := client.request(http.MethodGet, endpoint, nil, &result)

	return result, err
}

func (client *Client) ContentRepositoryCreate(hubID string, input ContentRepositoryInput) (ContentRepository, error) {
	result := ContentRepository{}
	body, err := json.Marshal(input)
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/hubs/%s/content-repositories", hubID)
	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentRepositoryUpdate(current ContentRepository, input ContentRepositoryInput) (ContentRepository, error) {
	result := ContentRepository{}

	body, err := createUpdatePatch(
		ContentRepositoryInput{
			Name:  current.Name,
			Label: current.Label,
		},
		input)

	if body == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/content-repositories/%s", current.ID)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentRepositoryList(hubID string, parameters PaginationParameters) (ContentRepositoryResults, error) {
	result := ContentRepositoryResults{}
	endpoint := fmt.Sprintf("/hubs/%s/content-repositories?%s", hubID, PaginationQueryString(parameters))
	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentRepositoryGetAll(hubID string) ([]ContentRepository, error) {
	parameters := PaginationParameters{}

	response, err := client.ContentRepositoryList(hubID, parameters)

	var result []ContentRepository
	result = append(result, response.Items...)

	for parameters.Page < response.Page.TotalPages-1 {
		parameters.Page++
		response, err := client.ContentRepositoryList(hubID, parameters)
		if err != nil {
			break
		}
		result = append(result, response.Items...)
	}

	return result, err
}

func (client *Client) ContentRepositoryFind() {

}

func (client *Client) ContentRepositoryShare() {

}

// ContentRepositoryAssignContentType assigns a Content Type to a Content Repository
func (client *Client) ContentRepositoryAssignContentType(repositoryID string, typeID string) (ContentRepository, error) {
	result := ContentRepository{}
	body, err := json.Marshal(struct {
		TypeID string `json:"contentTypeId"`
	}{
		typeID,
	})
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/content-repositories/%s/content-types", repositoryID)
	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

// ContentRepositoryRemoveContentType removes a Content Type from a Content Repository
func (client *Client) ContentRepositoryRemoveContentType(repositoryID string, typeID string) (ContentRepository, error) {
	result := ContentRepository{}
	endpoint := fmt.Sprintf("/content-repositories/%s/content-types/%s", repositoryID, typeID)
	err := client.request(http.MethodDelete, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentRepositoryAssignFeature(repositoryID string) {

}

func (client *Client) ContentRepositoryRemoveFeature(repositoryID string) {

}
