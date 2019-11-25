package internal

import (
	"testing"

	"github.com/operator-framework/api/pkg/validation/errors"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

func TestValidatePackageManifest(t *testing.T) {
	pkgName := "test-package"

	cases := []struct {
		validatorFuncTest
		pkg *registry.PackageManifest
	}{
		{
			validatorFuncTest{
				description: "successful validation",
			},
			&registry.PackageManifest{
				Channels: []registry.PackageChannel{
					{Name: "foo", CurrentCSVName: "bar"},
				},
				DefaultChannelName: "foo",
				PackageName:        "test-package",
			},
		},
		{
			validatorFuncTest{
				description: "successful validation no default channel with only one channel",
			},
			&registry.PackageManifest{
				Channels: []registry.PackageChannel{
					{Name: "foo", CurrentCSVName: "bar"},
				},
				PackageName: "test-package",
			},
		},
		{
			validatorFuncTest{
				description: "no default channel and more than one channel",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidPackageManifest("default channel is empty but more than one channel exists", pkgName),
				},
			},
			&registry.PackageManifest{
				Channels: []registry.PackageChannel{
					{Name: "foo", CurrentCSVName: "bar"},
					{Name: "foo2", CurrentCSVName: "baz"},
				},
				PackageName: "test-package",
			},
		},
		{
			validatorFuncTest{
				description: "default channel does not exist in channels",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidPackageManifest(`default channel "baz" not found in the list of declared channels`, pkgName),
				},
			},
			&registry.PackageManifest{
				Channels: []registry.PackageChannel{
					{Name: "foo", CurrentCSVName: "bar"},
				},
				DefaultChannelName: "baz",
				PackageName:        "test-package",
			},
		},
		{
			validatorFuncTest{
				description: "channels are empty",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidPackageManifest("channels empty", pkgName),
				},
			},
			&registry.PackageManifest{
				Channels:           nil,
				DefaultChannelName: "baz",
				PackageName:        "test-package",
			},
		},
		{
			validatorFuncTest{
				description: "one channel's CSVName is empty",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidPackageManifest(`channel "foo" currentCSV is empty`, pkgName),
				},
			},
			&registry.PackageManifest{
				Channels:           []registry.PackageChannel{{Name: "foo"}},
				DefaultChannelName: "foo",
				PackageName:        "test-package",
			},
		},
		{
			validatorFuncTest{
				description: "duplicate channel name",
				wantErr:     true,
				errors: []errors.Error{
					errors.ErrInvalidPackageManifest(`duplicate package manifest channel name "foo"`, pkgName),
				},
			},
			&registry.PackageManifest{
				Channels: []registry.PackageChannel{
					{Name: "foo", CurrentCSVName: "bar"},
					{Name: "foo", CurrentCSVName: "baz"},
				},
				DefaultChannelName: "foo",
				PackageName:        "test-package",
			},
		},
	}

	for _, c := range cases {
		result := validatePackageManifest(c.pkg)
		c.check(t, result)
	}
}
