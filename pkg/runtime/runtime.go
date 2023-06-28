package runtime

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/common"
	"github.com/imuxin/ksql/pkg/ext/abs"
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/pretty"
)

var (
	DefaultDatabase = ""
)

type Runnable[T any] interface {
	Run() ([]T, []pretty.PrintColumn, error)
}

var _ Runnable[any] = &RunnableImpl[any]{}

type RunnableImpl[T any] struct {
	ksql        *parser.KSQL
	restConfig  *rest.Config
	plugin      abs.Plugin
	whereFilter Filter
}

func NewDefaultRunnable[T any](ksql *parser.KSQL, restConfig *rest.Config) Runnable[T] {
	return &RunnableImpl[T]{
		ksql:       ksql,
		restConfig: restConfig,
	}
}

func (r RunnableImpl[T]) Run() (results []T, columns []pretty.PrintColumn, err error) {
	// TODO: 支持自定义 Table 拓展
	switch {
	case r.ksql.Desc != nil:
		return r.Desc()
	case r.ksql.Select != nil:
		return r.List()
	case r.ksql.Delete != nil:
		return r.Delete()
	default:
		return nil, nil, common.Unsupported()
	}
}
