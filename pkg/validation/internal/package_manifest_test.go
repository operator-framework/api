package internal

import (
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation/errors"
)

func TestValidatePackageManifest(t *testing.T) {
	pkgName := "test-package"

	cases := []struct {
		validatorFuncTest
		pkg *manifests.PackageManifest
	}{
		{
			validatorFuncTest{
				description: "successful validation",
			},
			&manifests.PackageManifest{
				Channels: []manifests.PackageChannel{
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
			&manifests.PackageManifest{
				Channels: []manifests.PackageChannel{
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
			&manifests.PackageManifest{
				Channels: []manifests.PackageChannel{
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
			&manifests.PackageManifest{
				Channels: []manifests.PackageChannel{
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
			&manifests.PackageManifest{
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
			&manifests.PackageManifest{
				Channels:           []manifests.PackageChannel{{Name: "foo"}},
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
			&manifests.PackageManifest{
				Channels: []manifests.PackageChannel{
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
