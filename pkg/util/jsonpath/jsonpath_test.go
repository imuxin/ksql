package jsonpath

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonPath(t *testing.T) {
	input := []byte(`{
		"kind": "Pod",
		"number": 1
	}`)
	var data interface{}
	assert.NoError(t, json.Unmarshal(input, &data))
	v, err := Find(data, "{.kind}")
	assert.NoError(t, err)
	assert.Equal(t, "Pod", v)

	v, err = Find(data, "{.number}")
	assert.NoError(t, err)
	assert.Equal(t, "1", v)

	data2 := struct {
		Kind   string `json:"kind"`
		Number int    `json:"number"`
	}{
		Kind:   "Pod",
		Number: 1,
	}

	v, err = Find(data2, "{.number}")
	assert.NoError(t, err)
	assert.Equal(t, "1", v)

}
