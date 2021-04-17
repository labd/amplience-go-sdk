package content

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
)

func createUpdatePatch(current interface{}, target interface{}) ([]byte, error) {

	src, err := json.Marshal(current)
	if err != nil {
		return nil, err
	}

	dst, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	if jsonpatch.Equal(src, dst) {
		return nil, nil
	}

	patch, err := jsonpatch.CreateMergePatch(src, dst)
	if err != nil {
		return nil, err
	}
	return patch, nil
}
