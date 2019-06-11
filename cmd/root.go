package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "operator-verify",
	Short: "New Manifest Verification Tool Prototype",
	Long:  `operator-verify is a CLI tool for the Operator Manifest Verification Library. This library provides functions to validate the operator manifest bundles against Operator-Lifecycle-Manager's ClusterServiceVersion type, CustomResourceDefinitions, and Package Manifest yamls. Currently, this application supports validation of ClusterServiceVersion yaml for any mismatched data types with Operator-Lifecycle-Manager's ClusterServiceVersion type.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing verification CLI tool...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
