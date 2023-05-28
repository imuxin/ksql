package repl

import (
	"fmt"
	"io"

	lop "github.com/samber/lo/parallel"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/imuxin/ksql/pkg/executor"
	"github.com/imuxin/ksql/pkg/pretty"
)

// Error ...
type Error string

// Errors
const (
	ErrContinue Error = "<continue input>"
	ErrQuit     Error = "<quit session>"
	ErrCmdRun   Error = "<command failed>"
)

func REPL() error {
	rl := newContLiner()
	defer rl.Close()
	for {
		in, err := rl.Prompt()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if in == "" {
			continue
		}

		if err := rl.Reindent(); err != nil {
			// fmt.Fprintf(g.errWriter, "error: %s\n", err)
			rl.Clear()
			continue
		}

		{
			result, err := executor.Execute[unstructured.Unstructured](in)
			if err != nil {
				fmt.Println(err)
			}

			r2 := lop.Map(result, func(item unstructured.Unstructured, index int) interface{} {
				return item.Object
			})

			pretty.Print(r2, []pretty.PrintColumn{
				{
					Name:     "NAME",
					JsonPath: "{ .metadata.name }",
				},
				{
					Name:     "NAMESPACE",
					JsonPath: "{ .metadata.namespace }",
				},
			})
		}

		// 	result, err := executor.Execute[appsv1.Deployment]("SELECT * FROM deploy NAMESPACE default")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// for _, item := range result {
		// 	repr.Println(item.Namespace, "/", item.Name)
		// }

		// err = s.Eval(in)
		// if err != nil {
		// 	if err == ErrContinue {
		// 		continue
		// 	} else if err == ErrQuit {
		// 		break
		// 	} else if err != ErrCmdRun {
		// 		rl.Clear()
		// 		continue
		// 	}
		// }
		rl.Accepted()
	}
	return nil
}
