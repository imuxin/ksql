package executor

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/compiler"
	"github.com/imuxin/ksql/pkg/pretty"
)

func ExecuteLikeSQL[T any](sql string, restConfig *rest.Config) ([]T, []pretty.PrintColumn, error) {
	runnable, err := compiler.Compile[T](sql, restConfig)
	if err != nil {
		return nil, nil, err
	}
	if runnable == nil {
		return []T{}, nil, nil
	}
	return runnable.Run()
}

func Execute[T any](sql string, restConfig *rest.Config) ([]T, error) {
	r, _, err := ExecuteLikeSQL[T](sql, restConfig)
	return r, err
}
