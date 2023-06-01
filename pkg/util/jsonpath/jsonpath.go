package jsonpath

import (
	"bytes"

	utiljsonpath "k8s.io/client-go/util/jsonpath"
)

func Find(data interface{}, pattern string) (string, error) {
	j := utiljsonpath.New("ksql-jsonpath")
	_ = j.AllowMissingKeys(true)
	if err := j.Parse(pattern); err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := j.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
