package compiler

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/runtime"
)

func Compile[T any](sql string, restConfig *rest.Config) (runtime.Runnable[T], error) {
	ksql, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if ksql == nil {
		return nil, nil
	}

	return runtime.NewDefaultRunnable[T](ksql, restConfig), nil
}
