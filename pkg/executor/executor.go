package executor

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/compiler"
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/pretty"
	"github.com/imuxin/ksql/pkg/runtime"
)

func Execute[T any](sql string, restConfig *rest.Config) ([]T, error) {
	runnable, err := compileToRunnable[T](sql, restConfig)
	if err != nil {
		return nil, err
	}
	if runnable == nil {
		return []T{}, nil
	}
	return runnable.Run()
}

func ExecuteLikeSQL[T any](sql string, restConfig *rest.Config) ([]pretty.PrintColumn, []T, error) {
	runable, err := compileToRunnable[T](sql, restConfig)
	if err != nil {
		return nil, nil, err
	}
	if runable == nil {
		return nil, []T{}, nil
	}
	return runable.RunLikeSQL()
}

func compileToRunnable[T any](sql string, restConfig *rest.Config) (runtime.Runnable[T], error) {
	ksql, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if ksql == nil {
		return nil, nil
	}
	return compiler.Compile[T](ksql, restConfig)
}
