package manifests

import (
	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation"
	"github.com/operator-framework/api/pkg/validation/errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "manifests <manifestDir>",
		Short: "Validate all manifests in a directory",
		Args:  cobra.ExactArgs(1),
		Long: `Validate the contents of a bundle in the supplied <manifestDir> directory.

For each manifest in the <manifestDir> directory, this command will print any errors
or warnings produced during the validation process.

Note: certain manifests may be ignored during the validation process if a validator
for that type/Kind has not been implemented yet in the Operator validation library.
`,
		Run: manifestsFunc,
	}

	rootCmd.Flags().Bool("operatorhub_validate", false, "enable optional UI validation for operatorhub.io")
	rootCmd.Flags().Bool("object_validate", false, "enable optional bundle object validation")

	return rootCmd
}

func manifestsFunc(cmd *cobra.Command, args []string) {
	bundle, err := manifests.GetBundleFromDir(args[0])
	if err != nil {
		log.Fatalf("Error generating bundle from directory: %s", err.Error())
	}
	if bundle == nil {
		log.Fatalf("Error generating bundle from directory")
	}

	operatorHubValidate, err := cmd.Flags().GetBool("operatorhub_validate")
	if err != nil {
		log.Fatalf("Unable to parse operatorhub_validate parameter: %v", err)
	}

	bundleObjectValidate, err := cmd.Flags().GetBool("object_validate")
	if err != nil {
		log.Fatalf("Unable to parse object_validate parameter: %v", err)
	}

	validators := validation.DefaultBundleValidators
	if operatorHubValidate {
		validators = validators.WithValidators(validation.OperatorHubValidator)
	}
	if bundleObjectValidate {
		validators = validators.WithValidators(validation.ObjectValidator)
	}

	results := validators.Validate(bundle.ObjectsToValidate()...)
	nonEmptyResults := []errors.ManifestResult{}
	for _, result := range results {
		if result.HasError() || result.HasWarn() {
			nonEmptyResults = append(nonEmptyResults, result)
		}
	}

	for _, result := range nonEmptyResults {
		for _, err := range result.Errors {
			log.Errorf(err.Error())
		}
		for _, err := range result.Warnings {
			log.Warnf(err.Error())
		}
	}

	return
}
