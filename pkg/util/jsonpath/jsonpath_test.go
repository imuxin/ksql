package jsonpath

import (
	"bytes"
	"testing"

	utiljsonpath "k8s.io/client-go/util/jsonpath"
)

func TestXxx(t *testing.T) {
	j := utiljsonpath.New("")
	j.AllowMissingKeys(true)
	buf := new(bytes.Buffer)
	err := j.Execute(buf, "")
	if err != nil {
		t.Error(err)
	}
}
