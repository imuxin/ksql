package ext

import (
	"strings"

	"github.com/imuxin/ksql/pkg/parser"
)

type Filter interface {
	Filter(i any) bool
}

type WhereFilters []*parser.Condition

var _ Filter = &WhereFilters{}

func (cs WhereFilters) Filter(i interface{}) bool {
	andList := make([]*parser.Condition, 0)
	orList := make([]*parser.Condition, 0)
	for _, item := range cs {
		switch strings.ToLower(item.Type) {
		case "and", "":
			andList = append(andList, item)
		case "or":
			orList = append(orList, item)
		}
	}

	for _, item := range orList {
		if item.Compare.Filter(i) {
			return true
		}
	}
	for _, item := range andList {
		if !item.Compare.Filter(i) {
			return false
		}
	}

	return true
}

func CompileWhereFilter(whereExpr *parser.WhereExpr) Filter {
	filter := make(WhereFilters, 0)
	if whereExpr == nil {
		return filter
	}
	filter = append(filter,
		&parser.Condition{
			Type:    "AND",
			Compare: whereExpr.First,
		},
	)
	filter = append(filter, whereExpr.Conditions...)
	return filter
}
