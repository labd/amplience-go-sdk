package content

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExtensionSnippet struct {
	Label string `json:"label,omitempty"`
	Body  string `json:"body,omitempty"`
}

type Extension struct {
	ID                        string             `json:"id,omitempty"`
	Name                      string             `json:"name"`
	Label                     string             `json:"label"`
	Description               string             `json:"description,omitempty"`
	URL                       string             `json:"url"`
	Height                    int                `json:"height,omitempty"`
	EnabledForAllContentTypes bool               `json:"enabledForAllContentTypes"`
	Category                  string             `json:"category"`
	Parameters                string             `json:"parameters,omitempty"`
	Snippets                  []ExtensionSnippet `json:"snippets,omitempty"`
	Settings                  string             `json:"settings,omitempty"`
	Links                     map[string]Link    `json:"_links,omitempty"`
}

type ExtensionInput struct {
	Name                      string             `json:"name"`
	Label                     string             `json:"label"`
	Description               string             `json:"description,omitempty"`
	URL                       string             `json:"url"`
	Height                    int                `json:"height,omitempty"`
	EnabledForAllContentTypes bool               `json:"enabledForAllContentTypes"`
	Category                  string             `json:"category"`
	Parameters                string             `json:"parameters,omitempty"`
	Snippets                  []ExtensionSnippet `json:"snippets,omitempty"`
	Settings                  string             `json:"settings,omitempty"`
}

type ExtensionResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []Extension
}

func (r *ExtensionResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["extensions"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (client *Client) ExtensionCreate(hubID string, input ExtensionInput) (Extension, error) {
	result := Extension{}
	body, err := json.Marshal(input)
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/hubs/%s/extensions", hubID)
	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

func (client *Client) ExtensionGet(id string) (Extension, error) {
	endpoint := fmt.Sprintf("/extensions/%s", id)
	result := Extension{}
	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ExtensionUpdate(current Extension, input ExtensionInput) (Extension, error) {
	result := Extension{}

	body, err := createUpdatePatch(
		ExtensionInput{
			Name:                      current.Name,
			Label:                     current.Label,
			Description:               current.Description,
			URL:                       current.URL,
			Height:                    current.Height,
			EnabledForAllContentTypes: current.EnabledForAllContentTypes,
			Category:                  current.Category,
			Parameters:                current.Parameters,
			Snippets:                  current.Snippets,
			Settings:                  current.Settings,
		},
		input)

	if body == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/extensions/%s", current.ID)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) ExtensionDelete(id string) error {
	endpoint := fmt.Sprintf("/extensions/%s", id)
	return client.request(http.MethodDelete, endpoint, nil, nil)
}

func (client *Client) ExtensionList(hubID string, parameters PaginationParameters) (ExtensionResults, error) {
	result := ExtensionResults{}
	endpoint := fmt.Sprintf("/hubs/%s/extensions?%s", hubID, PaginationQueryString(parameters))
	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ExtensionGetAll(hubID string) ([]Extension, error) {
	parameters := PaginationParameters{}
	response, err := client.ExtensionList(hubID, parameters)

	var result []Extension
	result = append(result, response.Items...)

	for parameters.Page < response.Page.TotalPages-1 {
		parameters.Page++
		response, err := client.ExtensionList(hubID, parameters)
		if err != nil {
			break
		}
		result = append(result, response.Items...)
	}

	return result, err
}
