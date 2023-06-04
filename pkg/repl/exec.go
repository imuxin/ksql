package repl

import (
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/executor"
	"github.com/imuxin/ksql/pkg/ext"
	"github.com/imuxin/ksql/pkg/pretty"
)

func Exec(in string, restConfig *rest.Config) (string, error) {
	columns, result, err := executor.ExecuteLikeSQL[ext.Object](in, restConfig)
	if err != nil {
		return "", err
	}
	return Format(columns, result), nil
}

func Format[T any](columns []pretty.PrintColumn, result []T) string {
	if len(result) == 0 {
		return "No rows to display"
	}

	return pretty.Format(result, columns)
}
