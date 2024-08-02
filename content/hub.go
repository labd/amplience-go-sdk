package content

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Hub struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name"`
	Label          string          `json:"label"`
	Description    *string         `json:"description,omitempty"`
	Plan           *string         `json:"plan,omitempty"`
	AlgoliaSearch  string          `json:"algoliaSearch"`
	CDV2           string          `json:"cdv2"`
	OrganizationID string          `json:"organizationId"`
	Settings       *Settings       `json:"settings,omitempty"`
	Links          map[string]Link `json:"_links"`
}

type AmplienceDamSettings struct {
	ApiKey    string `json:"API_KEY"`
	ApiSecret string `json:"API_SECRET"`
	Endpoint  string `json:"endpoint"`
}

type PlatformSettings struct {
	AmplienceDam *AmplienceDamSettings `json:"amplience_dam,omitempty"`
}

type PublishingSettings struct {
	Platforms *PlatformSettings `json:"platforms,omitempty"`
}

type DeviceSettings struct {
	Name      string `json:"name"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Orientate bool   `json:"orientate"`
}

type LocalizationSettings struct {
	Locales []string `json:"locales,omitempty"`
}

type ApplicationSettings struct {
	Name         string `json:"name"`
	TemplatedUri string `json:"templatedUri"`
}

type PreviewVirtualStagingEnvironmentSettings struct {
	Hostname string `json:"hostname"`
}

type VirtualStagingEnvironmentSettings struct {
	Hostname string `json:"hostname"`
}

type AssetManagementSettings struct {
	Enabled      *bool   `json:"enabled,omitempty"`
	ClientConfig *string `json:"clientConfig,omitempty"`
}

type Settings struct {
	Publishing                       *PublishingSettings                       `json:"publishing,omitempty"`
	Devices                          []DeviceSettings                          `json:"devices,omitempty"`
	Localization                     *LocalizationSettings                     `json:"localization,omitempty"`
	Applications                     []ApplicationSettings                     `json:"applications,omitempty"`
	PreviewVirtualStagingEnvironment *PreviewVirtualStagingEnvironmentSettings `json:"previewVirtualStagingEnvironment,omitempty"`
	VirtualStagingEnvironment        *VirtualStagingEnvironmentSettings        `json:"virtualStagingEnvironment,omitempty"`
	AssetManagement                  *AssetManagementSettings                  `json:"assetManagement,omitempty"`
}

type HubUpdateInput struct {
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	Description *string   `json:"description"`
	Settings    *Settings `json:"settings"`
}

type HubResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []Hub
}

func (r *HubResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["hubs"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (client *Client) HubList(parameters PaginationParameters) (HubResults, error) {
	parameters.Sort = "" // Sort is not supported.
	result := HubResults{}
	endpoint := fmt.Sprintf("/hubs?%s", PaginationQueryString(parameters))
	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) HubGetAll() ([]Hub, error) {
	parameters := PaginationParameters{}

	response, err := client.HubList(parameters)

	var result []Hub
	result = append(result, response.Items...)

	for parameters.Page < response.Page.TotalPages-1 {
		parameters.Page++
		response, err := client.HubList(parameters)
		if err != nil {
			break
		}
		result = append(result, response.Items...)
	}

	return result, err
}

// HubPatch will update hub settings. Note that if any settings are not provided they will be ignored during the
// patch, so they will continue existing.
func (client *Client) HubPatch(id string, input HubUpdateInput) (Hub, error) {
	endpoint := fmt.Sprintf("/hubs/%s", id)
	result := Hub{}

	body, err := json.Marshal(input)
	if err != nil {
		return Hub{}, err
	}

	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) HubGet(id string) (Hub, error) {
	endpoint := fmt.Sprintf("/hubs/%s", id)
	result := Hub{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}
