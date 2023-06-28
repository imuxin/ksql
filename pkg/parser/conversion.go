package parser

import (
	lop "github.com/samber/lo/parallel"

	"github.com/imuxin/ksql/pkg/pretty"
)

func (ksql *KSQL) CompilePrintColumns() []pretty.PrintColumn {
	if ksql.Select == nil {
		return nil
	}

	s := ksql.Select.Select
	if !s.ALL {
		return lop.Map(s.Columns, func(item *Column, index int) pretty.PrintColumn {
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
