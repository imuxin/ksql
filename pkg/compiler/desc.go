package compiler

import (
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
	// TODO
	return nil, nil
}
