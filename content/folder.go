package content

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Folder struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Links map[string]Link `json:"_links"`
}

type FolderInput struct {
	Name string `json:"name"`
}

// ContentItemResults is returned by the ContentItemList func
type FolderResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []Folder
}

// UnmarshalJSON is a custom unmarshaller for the embedded content
func (r *FolderResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["folders"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (client *Client) FolderCreate(repositoryID string, input FolderInput) (FolderInput, error) {
	result := FolderInput{}
	body, err := json.Marshal(input)
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/content-repositories/%s/folders", repositoryID)
	err = client.request(http.MethodPost, endpoint, body, &result)
	return result, err
}

func (client *Client) FolderGet(id string) (Folder, error) {
	endpoint := fmt.Sprintf("/folders/%s", id)
	result := Folder{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) FolderDelete(id string) (Folder, error) {
	endpoint := fmt.Sprintf("/folders/%s", id)
	result := Folder{}

	err := client.request(http.MethodDelete, endpoint, nil, &result)
	return result, err
}

func (client *Client) FolderList(repositoryID string) (FolderResults, error) {
	result := FolderResults{}
	endpoint := fmt.Sprintf("/content-repositories/%s/folders", repositoryID)

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}
