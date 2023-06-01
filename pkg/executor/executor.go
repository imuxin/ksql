package executor

import (
	"github.com/imuxin/ksql/pkg/compiler"
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/pretty"

	"k8s.io/client-go/rest"
)

func Execute[T any](sql string, restConfig *rest.Config) ([]T, error) {
	ksql, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if ksql == nil {
		return nil, nil
	}
	runable, err := compiler.Compile[T](ksql, restConfig)
	if err != nil {
		return nil, err
	}
	if runable == nil {
		return []T{}, nil
	}
	return runable.Run()
}

func ExecuteLikeSQL[T any](sql string, restConfig *rest.Config) ([]pretty.PrintColumn, []T, error) {
	ksql, err := parser.Parse(sql)
	if err != nil {
		return nil, nil, err
	}
	if ksql == nil {
		return nil, nil, nil
	}
	runable, err := compiler.Compile[T](ksql, restConfig)
	if err != nil {
		return nil, nil, err
	}
	if runable == nil {
		return nil, []T{}, nil
	}
	return runable.RunLikeSQL()
}
