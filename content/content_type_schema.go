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

type ContentTypeSchemaInput struct {
	SchemaID        string `json:"schemaId,omitempty"`
	Body            string `json:"body,omitempty"`
	ValidationLevel string `json:"validationLevel,omitempty"`
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

func (client *Client) ContentTypeSchemaCreate(hubID string, update ContentTypeSchemaInput) (ContentTypeSchema, error) {
	result := ContentTypeSchema{}
	body, err := json.Marshal(update)
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/hubs/%s/content-type-schemas", hubID)
	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentTypeSchemaGet(id string) (ContentTypeSchema, error) {
	endpoint := fmt.Sprintf("/content-type-schemas/%s", id)
	result := ContentTypeSchema{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentTypeSchemaUpdate(current ContentTypeSchema, update ContentTypeSchemaInput) (ContentTypeSchema, error) {
	result := ContentTypeSchema{}

	body, err := createUpdatePatch(
		ContentTypeSchemaInput{
			Body:            current.Body,
			ValidationLevel: current.ValidationLevel,
			SchemaID:        current.SchemaID,
		},
		update)

	if body == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/content-type-schemas/%s", current.ID)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentTypeSchemaList(hubID string) (ContentTypeSchemaResults, error) {
	result := ContentTypeSchemaResults{}
	endpoint := fmt.Sprintf("/hubs/%s/content-type-schemas", hubID)

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) ContentTypeSchemaArchive(id string, version int) (ContentTypeSchema, error) {
	result := ContentTypeSchema{}
	endpoint := fmt.Sprintf("/content-type-schemas/%s/archive", id)

	body, err := json.Marshal(ArchiveInput{Version: version})
	if err != nil {
		return result, err
	}

	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

func (client *Client) ContentTypeSchemaUnarchive(id string, version int) (ContentTypeSchema, error) {
	result := ContentTypeSchema{}
	endpoint := fmt.Sprintf("/content-type-schemas/%s/unarchive", id)

	body, err := json.Marshal(ArchiveInput{Version: version})
	if err != nil {
		return result, err
	}

	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}
