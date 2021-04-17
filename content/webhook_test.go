package content

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterSerialize(t *testing.T) {
	data := WebhookInput{
		Label:    "unittest",
		Events:   []string{"dynamic-content.edition.published"},
		Handlers: []string{"http://example.org/foobar"},
		Active:   true,
		Secret:   "my-secret",
		Filters: []WebhookFilter{
			WebhookFilterEqual{
				JSONPath: "$.payload.id",
				Value:    "1234",
			},
			WebhookFilterIn{
				JSONPath: "$.payload.id",
				Values: []string{
					"1234",
					"bar",
				},
			},
		},
		Method: "POST",
	}

	b, err := json.Marshal(&data)
	assert.Nil(t, err)

	expected := []byte(`
	{
		"label":"unittest",
		"events":[
			"dynamic-content.edition.published"
		],
		"handlers":[
			"http://example.org/foobar"
		],
		"active":true,
		"notifications":null,
		"secret":"my-secret",
		"filters":[
			{
				"type":"equal",
				"arguments":[
					{
						"jsonPath":"$.payload.id"
					},
					{
						"value":"1234"
					}
				]
			},
			{
				"type":"in",
				"arguments":[
					{
						"jsonPath":"$.payload.id"
					},
					{
						"value":[
							"1234",
							"bar"
						]
					}
				]
			}
		],
		"method":"POST"
	}
	`)

	assertJSONEqual(t, expected, b)

	result := WebhookInput{}
	if err := json.Unmarshal(b, &result); err != nil {
		assert.Nil(t, err)
	}
}
