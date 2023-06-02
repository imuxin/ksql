package compiler

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/common"
	"github.com/imuxin/ksql/pkg/parser"
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
	case ksql.Desc != nil:
		compiler = DescCompiler[T]{
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
