package main

import (
	"fmt"
	"os"

	"github.com/imuxin/ksql/pkg/repl"
)

func main() {
	if err := repl.REPL(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
