package repl

import (
	lop "github.com/samber/lo/parallel"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/imuxin/ksql/pkg/pretty"
)

func Output(result []unstructured.Unstructured) string {
	if len(result) == 0 {
		return "No rows to display"
	}

	r2 := lop.Map(result, func(item unstructured.Unstructured, index int) interface{} {
		return item.Object
	})

	return pretty.Format(r2, nil)
}
