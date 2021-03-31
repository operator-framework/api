package schema

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"cuelang.org/go/cue/load"
)

const schemaPath = "../manifests"

func TestValidateConfig(t *testing.T) {
	var table = []struct {
		description string
		filename    string
		kind        string
		hasError    bool
		errString   string
	}{
		{
			description: "valid bundle config",
			filename:    "./testdata/valid/bundle.json",
			kind:        olmBundle,
			hasError:    false,
		},
		{
			description: "valid package config",
			filename:    "./testdata/valid/package.json",
			kind:        olmPackage,
			hasError:    false,
		},
		{
			description: "valid channel property config",
			filename:    "./testdata/valid/channel_property.json",
			kind:        olmChannel,
			hasError:    false,
		},
		{
			description: "valid gvk property config",
			filename:    "./testdata/valid/gvk_property.json",
			kind:        olmGVKProvided,
			hasError:    false,
		},
		{
			description: "valid gvk property config",
			filename:    "./testdata/valid/package_property.json",
			kind:        olmPackageProperty,
			hasError:    false,
		},
		{
			description: "valid package config",
			filename:    "./testdata/valid/package.json",
			kind:        olmPackage,
			hasError:    false,
		},
		{
			description: "valid skiprange property config",
			filename:    "./testdata/valid/skiprange_property.json",
			kind:        olmSkipRange,
			hasError:    false,
		},
		{
			description: "valid skips prpoperty config",
			filename:    "./testdata/valid/skips_property.json",
			kind:        olmSkips,
			hasError:    false,
		},
		{
			description: "invalid bundle config",
			filename:    "./testdata/invalid/bundle.json",
			kind:        olmBundle,
			hasError:    true,
			errString:   `#olmbundle.relatedImages.0.name: incomplete value !=""`,
		},
		{
			description: "invalid channel property config",
			filename:    "./testdata/invalid/channel_property.json",
			kind:        olmChannel,
			hasError:    true,
			errString:   `#olmchannel.value.name: conflicting values !="" and 1 (mismatched types string and int)`,
		},
		{
			description: "invalid gvk property config",
			filename:    "./testdata/invalid/gvk_property.json",
			kind:        olmGVKProvided,
			hasError:    true,
			errString:   `#olmgvkprovided.value.version: incomplete value !=""`,
		},
		{
			description: "invalid package property config",
			filename:    "./testdata/invalid/package_property.json",
			kind:        olmPackageProperty,
			hasError:    true,
			errString:   `#packageproperty.value.packageName: incomplete value !=""`,
		},
		{
			description: "invalid package config",
			filename:    "./testdata/invalid/package.json",
			kind:        olmPackage,
			hasError:    true,
			errString:   `#olmpackage.defaultChannel: incomplete value !=""`,
		},
		{
			description: "invalid skiprange property config",
			filename:    "./testdata/invalid/skiprange_property.json",
			kind:        olmSkipRange,
			hasError:    true,
			errString:   `#olmskipRange.value: conflicting values !="" and 1 (mismatched types string and int)`,
		},
		{
			description: "invalid skips prpoperty config",
			filename:    "./testdata/invalid/skips_property.json",
			kind:        olmSkips,
			hasError:    true,
			errString:   `#olmskips.value: incomplete value !=""`,
		},
		{
			description: "mismatch schema config",
			filename:    "./testdata/invalid/bundle.json",
			kind:        olmPackage,
			hasError:    true,
			errString:   `#olmpackage.schema: conflicting values "olm.bundle" and "olm.package"`,
		},
	}

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			logger := logrus.NewEntry(logrus.New())
			// load schema for config definitions
			instance := load.Instances([]string{"."}, &load.Config{
				Dir: schemaPath,
			})
			if len(instance) > 1 {
				t.Fatalf("multiple instance loading currently not supported: %s", schemaPath)
			}
			if len(instance) < 1 {
				t.Fatalf("no instances found: %s", schemaPath)
			}

			// Config validator
			configValidator := NewConfigValidator(instance[0], logger)
			// Read json file
			content, err := ioutil.ReadFile(tt.filename)
			require.NoError(t, err)

			// Validate json against schema
			err = configValidator.Validate(content, tt.kind)

			if tt.hasError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
