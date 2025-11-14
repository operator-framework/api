package release

import (
	"encoding/json"
	"testing"

	semver "github.com/blang/semver/v4"
	"github.com/stretchr/testify/require"
)

func TestOperatorReleaseMarshal(t *testing.T) {
	tests := []struct {
		name string
		in   OperatorRelease
		out  []byte
		err  error
	}{
		{
			name: "single-segment",
			in:   OperatorRelease{Release: []semver.PRVersion{mustNewPRVersion("1")}},
			out:  []byte(`"1"`),
		},
		{
			name: "two-segments",
			in:   OperatorRelease{Release: []semver.PRVersion{mustNewPRVersion("1"), mustNewPRVersion("0")}},
			out:  []byte(`"1.0"`),
		},
		{
			name: "multi-segment",
			in: OperatorRelease{Release: []semver.PRVersion{
				mustNewPRVersion("1"),
				mustNewPRVersion("2"),
				mustNewPRVersion("3"),
			}},
			out: []byte(`"1.2.3"`),
		},
		{
			name: "numeric-segments",
			in: OperatorRelease{Release: []semver.PRVersion{
				mustNewPRVersion("20240101"),
				mustNewPRVersion("12345"),
			}},
			out: []byte(`"20240101.12345"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := tt.in.MarshalJSON()
			require.Equal(t, tt.out, m, string(m))
			require.Equal(t, tt.err, err)
		})
	}
}

func TestOperatorReleaseUnmarshal(t *testing.T) {
	type TestStruct struct {
		Release OperatorRelease `json:"r"`
	}
	tests := []struct {
		name string
		in   []byte
		out  TestStruct
		err  error
	}{
		{
			name: "single-segment",
			in:   []byte(`{"r": "1"}`),
			out:  TestStruct{Release: OperatorRelease{Release: []semver.PRVersion{mustNewPRVersion("1")}}},
		},
		{
			name: "two-segments",
			in:   []byte(`{"r": "1.0"}`),
			out:  TestStruct{Release: OperatorRelease{Release: []semver.PRVersion{mustNewPRVersion("1"), mustNewPRVersion("0")}}},
		},
		{
			name: "multi-segment",
			in:   []byte(`{"r": "1.2.3"}`),
			out: TestStruct{Release: OperatorRelease{Release: []semver.PRVersion{
				mustNewPRVersion("1"),
				mustNewPRVersion("2"),
				mustNewPRVersion("3"),
			}}},
		},
		{
			name: "numeric-segments",
			in:   []byte(`{"r": "20240101.12345"}`),
			out: TestStruct{Release: OperatorRelease{Release: []semver.PRVersion{
				mustNewPRVersion("20240101"),
				mustNewPRVersion("12345"),
			}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TestStruct{}
			err := json.Unmarshal(tt.in, &s)
			require.Equal(t, tt.out, s)
			require.Equal(t, tt.err, err)
		})
	}
}

func mustNewPRVersion(s string) semver.PRVersion {
	v, err := semver.NewPRVersion(s)
	if err != nil {
		panic(err)
	}
	return v
}
