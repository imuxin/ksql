package jsonpath

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	utiljsonpath "k8s.io/client-go/util/jsonpath"
)

func TestJsonPath(t *testing.T) {
	input := []byte(`{
		"kind": "Pod"
	}`)
	var data interface{}
	assert.NoError(t, json.Unmarshal(input, &data))
	j := utiljsonpath.New("hello")
	j.AllowMissingKeys(true)
	assert.NoError(t, j.Parse("{.kind}"))
	buf := new(bytes.Buffer)
	assert.NoError(t, j.Execute(buf, data))
	assert.Equal(t, "Pod", buf.String())
}
