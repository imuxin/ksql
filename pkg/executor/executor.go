package executor

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/compiler"
	"github.com/imuxin/ksql/pkg/pretty"
)

func ExecuteLikeSQL[T any](sql string, restConfig *rest.Config) ([]pretty.PrintColumn, []T, error) {
	runnable, printColumns, err := compiler.Compile[T](sql, restConfig)
	if err != nil {
		return nil, nil, err
	}
	if runnable == nil {
		return nil, []T{}, nil
	}
	r, err := runnable.Run()
	return printColumns, r, err
}

func Execute[T any](sql string, restConfig *rest.Config) ([]T, error) {
	_, r, err := ExecuteLikeSQL[T](sql, restConfig)
	return r, err
}
