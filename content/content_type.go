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

func (client *Client) ContentTypeList(hubID string) (ContentTypeResults, error) {
	result := ContentTypeResults{}
	endpoint := fmt.Sprintf("/hubs/%s/content-types", hubID)

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentTypeArchive(id string) (ContentType, error) {
	result := ContentType{}
	endpoint := fmt.Sprintf("/content-types/%s/archive", id)

	err := client.request(http.MethodPost, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentTypeUnarchive(id string, version int) (ContentType, error) {
	result := ContentType{}
	endpoint := fmt.Sprintf("/content-types/%s/unarchive", id)

	body, err := json.Marshal(ArchiveInput{Version: version})
	if err != nil {
		return result, err
	}

	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}
