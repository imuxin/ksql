package compiler

import (
	lop "github.com/samber/lo/parallel"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/common"
	"github.com/imuxin/ksql/pkg/ext"
	extkube "github.com/imuxin/ksql/pkg/ext/kube"
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/pretty"
	"github.com/imuxin/ksql/pkg/runtime"
)

func Compile[T any](ksql *parser.KSQL, restConfig *rest.Config) (runtime.Runnable[T], error) {
	var compiler Compiler[T]
	switch {
	case ksql.Use != nil:
		return nil, common.Unsupported()
	case ksql.Select != nil:
		compiler = SelectCompiler[T]{
			ksql:       ksql,
			restConfig: restConfig,
		}
	default: // TODO: support delete, update
		return nil, common.Unsupported()
	}
	return compiler.Compile()
}

type Compiler[T any] interface {
	Compile() (runtime.Runnable[T], error)
}

type SelectCompiler[T any] struct {
	ksql       *parser.KSQL
	restConfig *rest.Config
}

func (c SelectCompiler[T]) compileDownloader() (ext.Downloader, error) {
	names := make([]string, 0)
	selector := labels.NewSelector()
	for _, item := range c.ksql.Select.KubernetesFilters {
		switch {
		case item.Label != nil:
			r, err := (LabelCompare)(*item.Label).IntoRequirement()
			if err != nil {
				return nil, err
			}
			selector = selector.Add(*r)
		case item.Name != nil:
			names = append(names, *item.Name)
		}
	}

	return extkube.APIServerDownloader{
		RestConfig: c.restConfig,
		Table:      c.ksql.Select.From.Table,
		Namespace:  c.ksql.Select.Namespace,
		Names:      names,
		Selector:   selector,
	}, nil
}

func (c SelectCompiler[T]) compilePrintColumns() []pretty.PrintColumn {
	s := c.ksql.Select.Select
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

func (c SelectCompiler[T]) compileWhereFilter() runtime.Filter {
	filter := make(WhereFilters, 0)
	if c.ksql.Select.Where == nil {
		return filter
	}
	filter = append(filter,
		&parser.Condition{
			Type:    "AND",
			Compare: c.ksql.Select.Where.First,
		},
	)
	filter = append(filter, c.ksql.Select.Where.Conditions...)
	return filter
}

func (c SelectCompiler[T]) Compile() (runtime.Runnable[T], error) {
	if c.ksql.Select.From.Table == "" {
		return nil, nil
	}

	d, err := c.compileDownloader()
	if err != nil {
		return nil, err
	}

	return runtime.RunnableImpl[T]{
		Downloader:   d,
		WhereFilter:  c.compileWhereFilter(),
		PrintColumns: c.compilePrintColumns(),
	}, nil
}
