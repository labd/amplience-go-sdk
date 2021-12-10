package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AlgoliaIndexResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []AlgoliaIndex
}

type AlgoliaIndex struct {
	ID               string          `json:"id,omitempty"`
	ParentID         string          `json:"parentId,omitempty"`
	Label            string          `json:"label"`
	Name             string          `json:"name"`
	Type             string          `json:"type"` // PRODUCTION, STAGING
	Suffix           string          `json:"suffix"`
	ReplicaCount     int             `json:"replicaCount"`
	CreatedDate      *time.Time      `json:"createdDate,omitempty"`
	LastModifiedDate *time.Time      `json:"lastModifiedDate,omitempty"`
	Links            map[string]Link `json:"_links"`
}

type AssignedContentTypeResults struct {
	Links map[string]Link `json:"_links"`
	Page  PageInformation `json:"page"`
	Items []AssignedContentType
}

type AssignedContentType struct {
	ID               string          `json:"id,omitempty"`
	ContentTypeUri   string          `json:"contentTypeUri"`
	CreatedDate      *time.Time      `json:"createdDate,omitempty"`
	LastModifiedDate *time.Time      `json:"lastModifiedDate,omitempty"`
	Links            map[string]Link `json:"_links"`
}

type AlgoliaIndexInput struct {
	Suffix               string                     `json:"suffix"`
	Label                string                     `json:"label"`
	Type                 string                     `json:"type"` // PRODUCTION, STAGING
	AssignedContentTypes []AssignedContentTypeInput `json:"assignedContentTypes"`
}

type AssignedContentTypeInput struct {
	ContentTypeUri string `json:"contentTypeUri"`
}

func (r *AssignedContentTypeResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["assigned-content-types"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (r *AlgoliaIndexResults) UnmarshalJSON(data []byte) error {
	generic := GenericListResults{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if err := decodeStruct(generic.Embedded["indexes"], &r.Items); err != nil {
		return err
	}

	r.Links = generic.Links
	r.Page = generic.Page
	return nil
}

func (client *Client) AlgoliaIndexCreate(hub_id string, input AlgoliaIndexInput) (AlgoliaIndex, error) {
	result := AlgoliaIndex{}
	body, err := json.Marshal(input)
	if err != nil {
		return result, err
	}
	endpoint := fmt.Sprintf("/algolia-search/%s/indexes", hub_id)
	err = client.request(http.MethodPost, endpoint, body, &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (client *Client) AlgoliaIndexUpdate(hubID string, current AlgoliaIndex, input AlgoliaIndexInput) (AlgoliaIndex, error) {
	result := AlgoliaIndex{}

	body, err := createUpdatePatch(
		AlgoliaIndexInput{},
		input)

	if body == nil {
		return current, nil
	}

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/algolia-search/%s/indexes/%s", hubID, current.ID)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) AlgoliaIndexGet(hub_id string, id string) (AlgoliaIndex, error) {
	endpoint := fmt.Sprintf("/algolia-search/%s/indexes/%s", hub_id, id)
	result := AlgoliaIndex{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) AlgoliaIndexDelete(hub_id string, id string) (AlgoliaIndex, error) {
	endpoint := fmt.Sprintf("/algolia-search/%s/indexes/%s", hub_id, id)
	result := AlgoliaIndex{}

	err := client.request(http.MethodDelete, endpoint, nil, &result)
	return result, err
}

func (client *Client) AlgoliaIndexList(hub_id string) (AlgoliaIndexResults, error) {
	result := AlgoliaIndexResults{}
	endpoint := fmt.Sprintf("/algolia-search/%s/indexes", hub_id)
	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) AlgoliaIndexSettingsGet(hub_id string, id string) (AlgoliaIndexSettings, error) {
	endpoint := fmt.Sprintf("/algolia-search/%s/indexes/%s/settings", hub_id, id)
	result := AlgoliaIndexSettings{}

	err := client.request(http.MethodGet, endpoint, nil, &result)
	return result, err
}

func (client *Client) AlgoliaIndexSettingsUpdate(hub_id string, id string, input AlgoliaIndexSettings) (AlgoliaIndexSettings, error) {
	result := AlgoliaIndexSettings{}

	body, err := createUpdatePatch(AlgoliaIndexSettings{}, input)

	if err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("/algolia-search/%s/indexes/%s/settings", hub_id, id)
	err = client.request(http.MethodPatch, endpoint, body, &result)
	return result, err
}

func (client *Client) AlgoliaIndexWebhooksGet(hub_id string, id string) ([]Webhook, error) {
	endpoint := fmt.Sprintf("/algolia-search/%s/indexes/%s/assigned-content-types", hub_id, id)
	assignedContentTypes := AssignedContentTypeResults{}
	err := client.request(http.MethodGet, endpoint, nil, &assignedContentTypes)
	result := make([]Webhook, len(assignedContentTypes.Items))

	for i, item := range assignedContentTypes.Items {
		err = client.request(http.MethodGet, item.Links["webhook"].Href, nil, &result[i])
		if err != nil {
			return result, err
		}
	}

	return result, err
}

// https://www.algolia.com/doc/api-reference/api-parameters/
type AlgoliaIndexSettings struct {
	// NumericAttributesToIndex interface{} `json:"numericAttributesToIndex"`
	MinWordSizefor1Typo     int      `json:"minWordSizefor1Typo"`
	MinWordSizefor2Typos    int      `json:"minWordSizefor2Typos"`
	HitsPerPage             int      `json:"hitsPerPage"`
	MaxValuesPerFacet       int      `json:"maxValuesPerFacet"`
	Version                 int      `json:"version"`
	SearchableAttributes    []string `json:"searchableAttributes"`
	AttributesToRetrieve    []string `json:"attributesToRetrieve"`
	UnretrievableAttributes []string `json:"unretrievableAttributes"`
	OptionalWords           []string `json:"optionalWords"`
	AttributesForFaceting   []string `json:"attributesForFaceting"`
	AttributesToSnippet     []string `json:"attributesToSnippet"`
	AttributesToHighlight   []string `json:"attributesToHighlight"`
	PaginationLimitedTo     int      `json:"paginationLimitedTo"`
	AttributeForDistinct    string   `json:"attributeForDistinct"`
	ExactOnSingleWordQuery  string   `json:"exactOnSingleWordQuery"`
	Ranking                 []string `json:"ranking"`
	CustomRanking           []string `json:"customRanking"`
	SeparatorsToIndex       string   `json:"separatorsToIndex"`
	RemoveWordsIfNoResults  string   `json:"removeWordsIfNoResults"`
	QueryType               string   `json:"queryType"`
	HighlightPreTag         string   `json:"highlightPreTag"`
	HighlightPostTag        string   `json:"highlightPostTag"`
	SnippetEllipsisText     string   `json:"snippetEllipsisText"`
	AlternativesAsExact     []string `json:"alternativesAsExact"`
}
