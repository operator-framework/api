package cmd

import (
	"fmt"

	"github.com/dweepgogia/new-manifest-verification/pkg/validate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(verifyCmd)
}

var verifyCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Validate YAML against OLM's CSV type.",
	Long:  `Verifies the yaml file against Operator-Lifecycle-Manager's ClusterServiceVersion type. Reports errors for any mismatched data types. Takes in one argument i.e. path to the yaml file. Version: 1.0`,
	Run:   verifyFunc,
}

func verifyFunc(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		fmt.Printf("command %s requires exactly one argument", cmd.CommandPath())
	}

	manifestDirectory := args[0]

	_ = validate.ValidateManifest(manifestDirectory)
}
