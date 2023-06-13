package repl

import (
	"fmt"
	"io"
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

		res, err := Exec(in, nil)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Println(res)
		}

		rl.Accepted()
	}
	return nil
}
