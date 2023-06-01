package compiler

import (
	"strings"

	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/runtime"
)

type WhereFilters []*parser.Condition

var _ runtime.Filter = &WhereFilters{}

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
		if (Compare)(item.Compare).Filter(i) {
			return true
		}
	}
	for _, item := range andList {
		if !(Compare)(item.Compare).Filter(i) {
			return false
		}
	}

	return true
}

type Compare parser.Compare

func (c Compare) Filter(_ interface{}) bool {
	// TODO
	return true
}
