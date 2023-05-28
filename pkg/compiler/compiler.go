package compiler

import (
	"github.com/imuxin/ksql/pkg/parser"
	"github.com/imuxin/ksql/pkg/runtime"
)

func Compile[T any](ksql *parser.KSQL) runtime.Runnable[T] {
	d := runtime.APIServerDownloader{
		Table:     ksql.Select.From.Table,
		Namespace: ksql.Select.Namespace,
		Names:     ksql.Select.Name,
	}
	return runtime.KubernetesRunnable[T]{
		Downloader: d,
		Filter:     runtime.JsonPathFilter{},
	}
}
