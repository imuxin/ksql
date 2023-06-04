package main

import (
	"fmt"
	"os"

	"k8s.io/component-base/cli"

	"github.com/imuxin/ksql/pkg/cmd"
)

// func init() {
// 	defer klog.Flush() // flushes all pending log I/O

// 	klog.InitFlags(nil) // initializing the flags
// 	flag.Parse()        // parses the command-line flags
// 	// klog.Info("now you can see me")
// }

func main() {
	command := cmd.NewDefaultKSQLCommand()
	if err := cli.RunNoErrOutput(command); err != nil {
		// Pretty-print the error and exit with an error.
		fmt.Println(err)
		os.Exit(1)
	}
}
