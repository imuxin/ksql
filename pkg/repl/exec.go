package repl

import (
	"github.com/imuxin/ksql/pkg/executor"
	"github.com/imuxin/ksql/pkg/ext"
	"github.com/imuxin/ksql/pkg/pretty"
)

func Exec(in string) string {
	columns, result, err := executor.ExecuteLikeSQL[ext.Object](in, nil)
	if err != nil {
		return err.Error()
	}
	return Format(columns, result)
}

func Format[T any](columns []pretty.PrintColumn, result []T) string {
	if len(result) == 0 {
		return "No rows to display"
	}

	return pretty.Format(result, columns)
}
