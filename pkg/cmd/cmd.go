package cmd

import (
	"fmt"
	"os"

	"github.com/imuxin/ksql/pkg/repl"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"

	"github.com/imuxin/ksql/pkg/ext/kube"
)

type KSQLOptions struct {
	Arguments   []string
	ConfigFlags *genericclioptions.ConfigFlags

	genericclioptions.IOStreams
}

func NewDefaultKSQLCommand() *cobra.Command {
	return NewKSQLCommand(KSQLOptions{
		Arguments:   os.Args,
		ConfigFlags: kube.DefaultConfigFlags,
		IOStreams:   genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

func NewKSQLCommand(o KSQLOptions) *cobra.Command {
	warningHandler := rest.NewWarningWriter(o.IOStreams.ErrOut, rest.WarningWriterOptions{Deduplicate: true, Color: true})
	warningsAsErrors := false
	cmds := &cobra.Command{
		Use:   "ksql",
		Short: "ksql, a SQL-like language tool for kubernetes",
		Long: `
		ksql, a SQL-like language tool for kubernetes.

      Find more information at:
            https://github.com/imuxin/ksql`,
		Run: func(cmd *cobra.Command, args []string) {
			if Execute != "" {
				fmt.Println(repl.Exec(Execute, nil))
				return
			}
			if err := repl.REPL(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		// Hook before and after Run initialize and write profiles to disk,
		// respectively.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			rest.SetDefaultWarningHandler(warningHandler)
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			if err := flushProfiling(); err != nil {
				return err
			}
			if warningsAsErrors {
				count := warningHandler.WarningCount()
				switch count {
				case 0:
					// no warnings
				case 1:
					return fmt.Errorf("%d warning received", count)
				default:
					return fmt.Errorf("%d warnings received", count)
				}
			}
			return nil
		},
	}

	flags := cmds.PersistentFlags()
	addProfilingFlags(flags)
	flags.BoolVar(&warningsAsErrors, "warnings-as-errors", warningsAsErrors, "Treat warnings received from the server as errors and exit with a non-zero exit code")
	flags.StringVarP(&Execute, "execute", "e", "", "Execute the statement and quit")
	o.ConfigFlags.AddFlags(flags)
	return cmds
}
