package compiler

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/runtime"
)

var _ Compiler[any] = &DescCompiler[any]{}

type DescCompiler[T any] struct {
	ksql       *parser.KSQL
	restConfig *rest.Config
}

func (c DescCompiler[T]) Compile() (runtime.Runnable[T], error) {
	sql := fmt.Sprintf("SELECT * FROM crd NAME %s", c.ksql.Desc.Table)
	fmt.Println(sql)
	runnable, err := Compile[apiextensionsv1.CustomResourceDefinition](sql, c.restConfig)
	if err != nil {
		return nil, err
	}
	list, err := runnable.Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &runtime.DESCRunnableImpl[T]{
		Tables: list,
	}, nil
}
