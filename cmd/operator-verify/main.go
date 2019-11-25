package main

import (
	"fmt"
	"os"

	manifests "github.com/operator-framework/api/cmd/operator-verify/manifests"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "operator-verify",
		Short: "Operator manifest validation tool",
		Long:  `operator-verify is a CLI tool that calls functions in pkg/validation.`,
	}

	rootCmd.AddCommand(manifests.NewCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
