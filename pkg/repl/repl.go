package repl

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/imuxin/ksql/pkg/executor"
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
			fmt.Printf("error: %s\n", err)
			rl.Clear()
			continue
		}

		{
			result, err := executor.Execute[unstructured.Unstructured](in, nil)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(Output(result))
		}

		rl.Accepted()
	}
	return nil
}
