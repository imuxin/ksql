package executor

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/compiler"
	"github.com/imuxin/ksql/pkg/pretty"
)

func Execute[T any](sql string, restConfig *rest.Config) ([]T, error) {
	runnable, err := compiler.Compile[T](sql, restConfig)
	if err != nil {
		return nil, err
	}
	if runnable == nil {
		return []T{}, nil
	}
	return runnable.Run()
}

func ExecuteLikeSQL[T any](sql string, restConfig *rest.Config) ([]pretty.PrintColumn, []T, error) {
	runable, err := compiler.Compile[T](sql, restConfig)
	if err != nil {
		return nil, nil, err
	}
	if runable == nil {
		return nil, []T{}, nil
	}
	return runable.RunLikeSQL()
}
