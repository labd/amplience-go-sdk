package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ContentTypeSchema struct {
	ID               string          `json:"id"`
	Body             string          `json:"body"`
	Version          int             `json:"version"`
	Status           string          `json:"status"`
	CreatedBy        string          `json:"createdBy"`
	CreatedDate      *time.Time      `json:"createdDate"`
	LastModifiedBy   string          `json:"lastModifiedBy"`
	LastModifiedDate *time.Time      `json:"lastModifiedDate"`
	Links            map[string]Link `json:"_links"`
	SchemaID         string          `json:"schemaId"`
	ValidationLevel  string          `json:"validationLevel"`
}

type ContentTypeSchemaUpdate struct {
	Body            string `json:"body"`
	ValidationLevel string `json:"validationLevel"`
}

type ContentTypeSchemaResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []ContentTypeSchema
}

func (r *ContentTypeSchemaResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["content-type-schemas"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (client *Client) ContentTypeSchemaGet(id string) (ContentTypeSchema, error) {
	endpoint := fmt.Sprintf("/content-type-schemas/%s", id)
	result := ContentTypeSchema{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err

}

func (client *Client) ContentTypeSchemaUpdate(current ContentTypeSchema, update ContentTypeSchemaUpdate) (ContentTypeSchema, error) {
	result := ContentTypeSchema{}

	patchBody, err := createUpdatePatch(
		ContentTypeSchemaUpdate{
			Body:            current.Body,
			ValidationLevel: current.ValidationLevel,
		},
		update)

	if patchBody == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/content-type-schemas/%s", current.ID)

	err = client.request(http.MethodPatch, endpoint, patchBody, &result)
	return result, err

}

func (client *Client) ContentTypeSchemaCreate() {

}

func (client *Client) ContentTypeSchemaList(hubId string) (ContentTypeSchemaResults, error) {
	result := ContentTypeSchemaResults{}
	endpoint := fmt.Sprintf("/hubs/%s/content-type-schemas", hubId)

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}
