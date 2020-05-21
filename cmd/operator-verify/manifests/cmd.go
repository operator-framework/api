package manifests

import (
	"fmt"
	"os"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation"
	"github.com/operator-framework/api/pkg/validation/errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "manifests",
		Short: "Validates all manifests in a directory",
		Long: `'operator-verify manifests' validates all manifests in the supplied directory
and prints errors and warnings corresponding to each manifest found to be
invalid. Manifests are only validated if a validator for that manifest
type/kind, ex. CustomResourceDefinition, is implemented in the Operator
validation library.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("command %s requires exactly one argument", cmd.CommandPath())
			}
			bundle, err := manifests.GetBundleFromDir(args[0])
			if err != nil {
				log.Fatalf("Error generating bundle from directory %s", err.Error())
			}
			results := validation.AllValidators.Validate(bundle)
			nonEmptyResults := []errors.ManifestResult{}
			for _, result := range results {
				if result.HasError() || result.HasWarn() {
					nonEmptyResults = append(nonEmptyResults, result)
				}
			}
			if len(nonEmptyResults) != 0 {
				fmt.Println(nonEmptyResults)
				os.Exit(1)
			}
		},
	}
}
