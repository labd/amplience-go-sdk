package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookEventsEnum string

const (
	WebhookContentItemAssigned        WebhookEventsEnum = "dynamic-content.content-item.assigned"
	WebhookContentItemCreated         WebhookEventsEnum = "dynamic-content.content-item.created"
	WebhookContentItemUpdated         WebhookEventsEnum = "dynamic-content.content-item.updated"
	WebhookContentItemWorkflowUpdated WebhookEventsEnum = "dynamic-content.content-item.workflow.updated"
	WebhookEditionPublished           WebhookEventsEnum = "dynamic-content.edition.published"
	WebhookEditionScheduled           WebhookEventsEnum = "dynamic-content.edition.scheduled"
	WebhookEditionUnscheduled         WebhookEventsEnum = "dynamic-content.edition.unscheduled"
	WebhookSnapshotPublished          WebhookEventsEnum = "dynamic-content.snapshot.published"
)

type Webhook struct {
	ID               string                `json:"id,omitempty"`
	Label            string                `json:"label"`
	Events           []string              `json:"events"`
	Handlers         []string              `json:"handlers"`
	Active           bool                  `json:"active"`
	Notifications    []Notification        `json:"notifications"`
	Secret           string                `json:"secret"`
	CreatedDate      *time.Time            `json:"createdDate,omitempty"`
	LastModifiedDate *time.Time            `json:"lastModifiedDate,omitempty"`
	Headers          []WebhookHeader       `json:"headers,omitempty"`
	Filters          []WebhookFilter       `json:"filters,omitempty"`
	Method           string                `json:"method"`
	CustomPayload    *WebhookCustomPayload `json:"customPayload,omitempty"`
}

func (obj *Webhook) UnmarshalJSON(data []byte) error {
	type Alias Webhook
	if err := json.Unmarshal(data, (*Alias)(obj)); err != nil {
		return err
	}

	for i := range obj.Filters {
		var err error
		obj.Filters[i], err = mapWebhookFilter(obj.Filters[i])
		if err != nil {
			return err
		}
	}
	return nil
}

type WebhookResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []Webhook
}

func (r *WebhookResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["webhooks"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

type WebhookInput struct {
	Label         string                `json:"label"`
	Events        []string              `json:"events"`
	Handlers      []string              `json:"handlers"`
	Active        bool                  `json:"active"`
	Notifications []Notification        `json:"notifications"`
	Secret        string                `json:"secret"`
	Headers       []WebhookHeader       `json:"headers,omitempty"`
	Filters       []WebhookFilter       `json:"filters,omitempty"`
	Method        string                `json:"method"`
	CustomPayload *WebhookCustomPayload `json:"customPayload,omitempty"`
}

func (obj *WebhookInput) UnmarshalJSON(data []byte) error {
	type Alias WebhookInput
	if err := json.Unmarshal(data, (*Alias)(obj)); err != nil {
		return err
	}

	for i := range obj.Filters {
		var err error
		obj.Filters[i], err = mapWebhookFilter(obj.Filters[i])
		if err != nil {
			return err
		}
	}
	return nil
}

type Notification struct {
	Email string
}

type WebhookHeader struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}

type WebhookCustomPayload struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (client *Client) WebhookCreate(hubID string, input WebhookInput) (Webhook, error) {
	endpoint := fmt.Sprintf("/hubs/%s/webhooks", hubID)
	result := Webhook{}

	body, err := json.Marshal(input)
	if err != nil {
		return result, err
	}

	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

func (client *Client) WebhookGet(hubID string, ID string) (Webhook, error) {
	endpoint := fmt.Sprintf("/hubs/%s/webhooks/%s", hubID, ID)
	result := Webhook{}
	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) WebhookUpdate(hubID string, current Webhook, input WebhookInput) (Webhook, error) {
	result := Webhook{}

	body, err := createUpdatePatch(
		WebhookInput{
			Label:         current.Label,
			Events:        current.Events,
			Handlers:      current.Handlers,
			Active:        current.Active,
			Secret:        current.Secret,
			Filters:       current.Filters,
			Method:        current.Method,
			Notifications: current.Notifications,
		},
		input)

	if body == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/hubs/%s/webhooks/%s", hubID, current.ID)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) WebhookDelete(hub_id string, id string) error {
	endpoint := fmt.Sprintf("/hubs/%s/webhooks/%s", hub_id, id)
	err := client.request(http.MethodDelete, endpoint, nil, nil)
	return err
}

func (client *Client) WebhookList(hub_id string, parameters PaginationParameters) (WebhookResults, error) {
	result := WebhookResults{}
	endpoint := fmt.Sprintf("/hubs/%s/webhooks?%s", hub_id, PaginationQueryString(parameters))
	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) WebhookGetAll(hubID string) ([]Webhook, error) {
	parameters := PaginationParameters{}

	response, err := client.WebhookList(hubID, parameters)

	var result []Webhook
	result = append(result, response.Items...)

	for parameters.Page < response.Page.TotalPages-1 {
		parameters.Page++
		response, err := client.WebhookList(hubID, parameters)
		if err != nil {
			break
		}
		result = append(result, response.Items...)
	}

	return result, err
}
