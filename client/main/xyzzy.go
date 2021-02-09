package main

import (
	"github.com/cyberax/go-dd-service-base/utils"
	"github.com/cyberax/go-dd-service-base/visibility"
	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
	"xyzzy/client"
)

func main() {
	// Stop the tracing
	tracer.Start(tracer.WithLogger(visibility.NopLogger{}))
	defer tracer.Stop()

	cobra.EnableCommandSorting = false

	var rootCmd = &cobra.Command{
		Use: "xyzzy",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.CheckRequiredFlags(cmd)
		},
	}

	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().String("server", "", "The Xyzzy server")
	rootCmd.PersistentFlags().String("token", os.Getenv("XYZZY_TOKEN"), "The token ID")
	_ = rootCmd.MarkPersistentFlagRequired("server")

	rootCmd.Flags().SortFlags = false

	rootCmd.AddCommand(utils.MakeCompletionCmd())
	rootCmd.AddCommand(client.MakePingCmd())

	err := rootCmd.Execute()
	if err != nil {
		tracer.Stop()
		_, _ = os.Stderr.Write([]byte(err.Error()))
		_, _ = os.Stderr.Write([]byte("\n"))
		os.Exit(1)
	}
}
