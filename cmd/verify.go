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
	Use:   "verify",
	Short: "Validate YAML against OLM's CSV type.",
	Long:  `Verifies the yaml file against Operator-Lifecycle-Manager's ClusterServiceVersion type. Reports errors for any mismatched data types. Takes in one argument i.e. path to the yaml file. Version: 1.0`,
	RunE:  verifyFunc,
}

func verifyFunc(cmd *cobra.Command, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("command %s requires exactly one argument", cmd.CommandPath())
	}

	yamlFileName := args[0]

	return validate.Verify(yamlFileName)
}
