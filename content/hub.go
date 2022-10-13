package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Hub struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Status      string `json:"status"`

	CreatedBy        string     `json:"createdBy"`
	CreatedDate      *time.Time `json:"createdDate"`
	LastModifiedBy   string     `json:"lastModifiedBy"`
	LastModifiedDate *time.Time `json:"lastModifiedDate"`
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

func (client *Client) HubCreate() {
}

func (client *Client) HubUpdate() {

}

func (client *Client) HubGet(ID string) (Hub, error) {
	endpoint := fmt.Sprintf("/hubs/%s", ID)
	result := Hub{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}
