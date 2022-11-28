package content

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ContentTypeIcon struct {
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type ContentTypeVisualization struct {
	Label        string `json:"label"`
	TemplatedURI string `json:"templatedUri"`
	Default      bool   `json:"default"`
}

type ContentTypeSettings struct {
	Label          string                     `json:"label"`
	Icons          []ContentTypeIcon          `json:"icons,omitempty"`
	Visualizations []ContentTypeVisualization `json:"visualizations,omitempty"`
}

type ContentType struct {
	ID             string              `json:"id"`
	ContentTypeURI string              `json:"contentTypeUri"`
	Status         string              `json:"status"`
	Settings       ContentTypeSettings `json:"settings,omitempty"`
	Links          map[string]Link     `json:"_links"`
}

type ContentTypeInput struct {
	ContentTypeURI string              `json:"contentTypeUri"`
	Settings       ContentTypeSettings `json:"settings,omitempty"`
}

type ContentTypeResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []ContentType
}

type ContentTypeSyncResult struct {
	ContentTypeURI string `json:"contentTypeUri"`
}

func (r *ContentTypeResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["content-types"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (client *Client) ContentTypeCreate(hubID string, input ContentTypeInput) (ContentType, error) {
	result := ContentType{}
	body, err := json.Marshal(input)
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/hubs/%s/content-types", hubID)
	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentTypeGet(id string) (ContentType, error) {
	endpoint := fmt.Sprintf("/content-types/%s", id)
	result := ContentType{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentTypeFindByUri(uri string, hubId string) (ContentType, error) {
	dummy := ContentType{}
	allItems, getErr := client.ContentTypeGetAll(hubId, StatusAny)

	if getErr != nil {
		return dummy, getErr
	}

	for _, item := range allItems {
		if item.ContentTypeURI == uri {
			return item, nil
		}
	}

	return dummy, fmt.Errorf(fmt.Sprintf("Could not find content-type %s", uri))
}

func (client *Client) ContentTypeUpdate(current ContentType, input ContentTypeInput) (ContentType, error) {
	result := ContentType{}

	body, err := createUpdatePatch(
		ContentTypeInput{
			Settings: current.Settings,
		},
		input)

	if body == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/content-types/%s", current.ID)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentTypeList(hubID string, parameters StatusPaginationParameters) (ContentTypeResults, error) {
	result := ContentTypeResults{}
	endpoint := fmt.Sprintf("/hubs/%s/content-types?%s", hubID, ContentTypePaginationQueryString(parameters))

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentTypeSyncSchema(current ContentType) (ContentTypeSyncResult, error) {
	result := ContentTypeSyncResult{}

	var emptyInput struct{}
	body, err := json.Marshal(emptyInput)

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/content-types/%s/schema", current.ID)

	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentTypeGetAll(hubID string, status ContentStatus) ([]ContentType, error) {
	parameters := StatusPaginationParameters{}
	parameters.Status = status

	response, err := client.ContentTypeList(hubID, parameters)

	var result []ContentType
	result = append(result, response.Items...)

	for parameters.Page < response.Page.TotalPages-1 {
		parameters.Page++
		response, err := client.ContentTypeList(hubID, parameters)
		if err != nil {
			break
		}
		result = append(result, response.Items...)
	}

	return result, err
}

func (client *Client) ContentTypeArchive(id string) (ContentType, error) {
	result := ContentType{}
	endpoint := fmt.Sprintf("/content-types/%s/archive", id)

	err := client.request(http.MethodPost, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentTypeUnarchive(id string) (ContentType, error) {
	result := ContentType{}
	endpoint := fmt.Sprintf("/content-types/%s/unarchive", id)

	err := client.request(http.MethodPost, endpoint, nil, &result)
	return result, err
}
