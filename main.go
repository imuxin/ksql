package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/imuxin/ksql/pkg/repl"
	"k8s.io/klog/v2"
)

func init() {
	defer klog.Flush() // flushes all pending log I/O

	klog.InitFlags(nil) // initializing the flags
	flag.Parse()        // parses the command-line flags
	// klog.Info("now you can see me")
}

func main() {
	if err := repl.REPL(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
