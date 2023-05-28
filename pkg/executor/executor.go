package executor

import (
	"github.com/imuxin/ksql/pkg/compiler"
	"github.com/imuxin/ksql/pkg/parser"
)

func Execute[T any](sql string) ([]T, error) {
	ksql, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if ksql == nil {
		return nil, nil
	}
	return compiler.Compile[T](ksql).Run()
}
