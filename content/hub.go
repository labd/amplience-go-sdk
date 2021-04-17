package content

import (
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

func (client *Client) HubList() {
	client.request(http.MethodGet, "/hubs", nil)
}

func (client *Client) HubCreate() {
}

func (client *Client) HubUpdate() {

}

func (client *Client) HubGet(ID string) (Hub, error) {
	endpoint := fmt.Sprintf("/hubs/%s", ID)
	result := Hub{}

	err := client.request(http.MethodGet, endpoint, &result)
	return result, err
}
