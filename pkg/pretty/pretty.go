package pretty

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	lop "github.com/samber/lo/parallel"

	"github.com/imuxin/ksql/pkg/util/jsonpath"
)

type PrintColumn struct {
	Name     string
	JsonPath string
}

func ToGenericArray(arr ...interface{}) []interface{} {
	return arr
}

func Print[T any](list []T, columns []PrintColumn) {
	t := table.NewWriter()
	headers := lop.Map(columns, func(item PrintColumn, index int) interface{} {
		return item.Name
	})
	t.AppendHeader(headers)

	for _, item := range list {
		raw := lop.Map(columns, func(_item PrintColumn, index int) interface{} {
			v, _ := jsonpath.Find(item, _item.JsonPath)
			return v
		})
		t.AppendRow(raw)
	}
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	fmt.Println(t.Render())
}
