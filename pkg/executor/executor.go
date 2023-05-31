package executor

import (
	"github.com/imuxin/ksql/pkg/compiler"
	"github.com/imuxin/ksql/pkg/parser"

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
	return runable.Run()
}
