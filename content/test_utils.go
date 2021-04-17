package content

import (
	"encoding/json"
	"reflect"

	"github.com/stretchr/testify/assert"
)

func assertJSONEqual(t assert.TestingT, a, b []byte) {
	var x, y interface{}
	var err error

	err = json.Unmarshal(a, &x)
	assert.NoError(t, err)

	err = json.Unmarshal(b, &y)
	assert.NoError(t, err)

	equal := reflect.DeepEqual(x, y)
	assert.True(t, equal)
}
