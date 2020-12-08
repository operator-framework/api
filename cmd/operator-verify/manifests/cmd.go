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
		Use:   "manifests",
		Short: "Validates all manifests in a directory",
		Long: `'operator-verify manifests' validates a bundle in the supplied directory
and prints errors and warnings corresponding to each manifest found to be
invalid. Manifests are only validated if a validator for that manifest
type/kind, ex. CustomResourceDefinition, is implemented in the Operator
validation library.`,
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
		log.Fatalf("Unable to parse operatorhub_validate parameter")
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
