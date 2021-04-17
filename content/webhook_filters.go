package content

import (
	"encoding/json"
	"errors"
)

type WebhookFilter interface{}

func mapWebhookFilter(input WebhookFilter) (WebhookFilter, error) {
	data, ok := input.(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid data")
	}

	discriminator, ok := data["type"].(string)
	if !ok {
		return nil, errors.New("Error processing discriminator field 'type'")
	}

	args, ok := data["arguments"].([]interface{})
	if !ok {
		return nil, errors.New("Error processing arguments")
	}

	switch discriminator {
	case "equal":
		new := WebhookFilterEqual{Type: "equal"}
		for _, arg := range args {
			mapArg := arg.(map[string]interface{})
			if value, ok := mapArg["jsonPath"]; ok {
				new.JSONPath = value.(string)
			}
			if value, ok := mapArg["value"]; ok {
				new.Value = value.(string)
			}
		}
		return new, nil

	case "in":
		new := WebhookFilterIn{Type: "in"}
		for _, arg := range args {
			mapArg := arg.(map[string]interface{})
			if value, ok := mapArg["jsonPath"]; ok {
				new.JSONPath = value.(string)
			}
			if value, ok := mapArg["value"].([]interface{}); ok {
				new.Values = make([]string, len(value))
				for i, valItem := range value {
					new.Values[i] = valItem.(string)
				}
			}
		}
		return new, nil
	}
	return nil, nil

}

type WebhookFilterEqual struct {
	Type     string `json:"type"`
	JSONPath string
	Value    string
}

func (obj WebhookFilterEqual) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string        `json:"type"`
		Arguments []interface{} `json:"arguments"`
	}{
		Type: "equal",
		Arguments: []interface{}{
			struct {
				JSONPath string `json:"jsonPath"`
			}{
				JSONPath: obj.JSONPath,
			},
			struct {
				Value string `json:"value"`
			}{
				Value: obj.Value,
			},
		},
	})
}

type WebhookFilterIn struct {
	Type     string `json:"type"`
	JSONPath string
	Values   []string `json:"values"`
}

func (obj WebhookFilterIn) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string        `json:"type"`
		Arguments []interface{} `json:"arguments"`
	}{
		Type: "in",
		Arguments: []interface{}{
			struct {
				JSONPath string `json:"jsonPath"`
			}{
				JSONPath: obj.JSONPath,
			},
			struct {
				Value []string `json:"value"`
			}{
				Value: obj.Values,
			},
		},
	})
}
