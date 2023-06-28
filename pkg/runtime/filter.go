package runtime

import (
	"strings"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/imuxin/ksql/pkg/parser"
)

type Filter interface {
	Filter(i any) bool
}

type WhereFilters []*parser.Condition

var _ Filter = &WhereFilters{}

func (f WhereFilters) Filter(i interface{}) bool {
	andList := make([]*parser.Condition, 0)
	orList := make([]*parser.Condition, 0)
	for _, item := range f {
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

type KubernetesFilters []*parser.KubernetesFilter

func (f KubernetesFilters) Convert() ([]string, labels.Selector, error) {
	names := make([]string, 0)
	selector := labels.NewSelector()
	for _, item := range f {
		switch {
		case item.Label != nil:
			req, err := item.Label.IntoRequirement()
			if err != nil {
				return nil, nil, err
			}
			selector = selector.Add(*req)
		case item.Name != nil:
			names = append(names, *item.Name)
		}
	}
	return names, selector, nil
}
