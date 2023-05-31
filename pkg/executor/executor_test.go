package executor

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestXxx(_ *testing.T) {
	sql := `select * from deploy namespace default label app = "zookeeper" label a in ("1", "2")`
	_, _ = Execute[unstructured.Unstructured](sql, nil)
}
