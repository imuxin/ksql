package compiler

import (
	lop "github.com/samber/lo/parallel"

	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/common"
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/pretty"
	"github.com/imuxin/ksql/pkg/runtime"
)

func Compile[T any](sql string, restConfig *rest.Config) (runtime.Runnable[T], []pretty.PrintColumn, error) {
	ksql, err := parser.Parse(sql)
	if err != nil {
		return nil, nil, err
	}
	if ksql == nil {
		return nil, nil, nil
	}

	r := runtime.NewDefaultRunnable[T](ksql, restConfig)
	var printColumns []pretty.PrintColumn
	switch {
	case ksql.Use != nil:
		return nil, nil, common.Unsupported()
	case ksql.Select != nil:
		printColumns = compilePrintColumns(ksql)
	case ksql.Delete != nil:
	case ksql.Desc != nil:
		printColumns = []pretty.PrintColumn{
			{
				Name:     "SCHEMA",
				JSONPath: "{ .spec }",
			},
			{
				Name:     "VERSION",
				JSONPath: "{ .version }",
			},
		}
	default: // TODO: support update
		return nil, nil, common.Unsupported()
	}
	return r, printColumns, err
}

func compilePrintColumns(ksql *parser.KSQL) []pretty.PrintColumn {
	s := ksql.Select.Select
	if !s.ALL {
		return lop.Map(s.Columns, func(item *parser.Column, index int) pretty.PrintColumn {
			name := item.Alias
			if name == "" {
				name = item.Name
			}
			return pretty.PrintColumn{
				Name:     name,
				JSONPath: item.Name,
			}
		})
	}
	return nil
}
