package pretty

import (
	"github.com/jedib0t/go-pretty/v6/table"
	lop "github.com/samber/lo/parallel"

	"github.com/imuxin/ksql/pkg/util/jsonpath"
)

type PrintColumn struct {
	Name     string
	JSONPath string
}

var defaultPrintColumns = []PrintColumn{
	{
		Name:     "NAME",
		JSONPath: "{ .metadata.name }",
	},
	{
		Name:     "NAMESPACE",
		JSONPath: "{ .metadata.namespace }",
	},
}

// func ToGenericArray(arr ...interface{}) []interface{} {
// 	return arr
// }

func Format[T any](list []T, columns []PrintColumn) string {
	if len(columns) == 0 {
		columns = defaultPrintColumns
	}

	t := table.NewWriter()
	headers := lop.Map(columns, func(item PrintColumn, index int) interface{} {
		return item.Name
	})
	t.AppendHeader(headers)

	for _, item := range list {
		raw := lop.Map(columns, func(_item PrintColumn, index int) interface{} {
			v, _ := jsonpath.Find(item, _item.JSONPath)
			return v
		})
		t.AppendRow(raw)
	}
	// t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	return t.Render()
}
