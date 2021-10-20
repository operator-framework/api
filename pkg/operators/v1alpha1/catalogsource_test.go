package v1alpha1

import (
	"encoding/json"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPollingIntervalDuration(t *testing.T) {
	var validDuration, _ = time.ParseDuration("45m")
	var defaultDuration, _ = time.ParseDuration("15m")

	type TestStruct struct {
		UpdateStrategy *UpdateStrategy `json:"updateStrategy,omitempty"`
	}

	tests := []struct {
		in  []byte
		out *UpdateStrategy
		err error
		empty bool
	}{
		{
			in: []byte(`{"UpdateStrategy": {"registryPoll":{"interval":"45m"}}}`),
			out: &UpdateStrategy{
				RegistryPoll: &RegistryPoll{
					Interval: &metav1.Duration{Duration: validDuration},
				},
			},
			err: nil,
		},
		{
			in: []byte(`{"UpdateStrategy": {"registryPoll":{"interval":"10m Error"}}}`),
			out: &UpdateStrategy{
				RegistryPoll: &RegistryPoll{
					Interval: &metav1.Duration{Duration: defaultDuration},
				},
			},
			err: nil,
		},
		{
			in:  []byte(`{"UpdateStrategy": {}}`),
			err: nil,
			empty: true,
		},
	}

	for _, tt := range tests {
		tc := TestStruct{}
		err := json.Unmarshal(tt.in, &tc)
		if err != tt.err {
			t.Fatalf("during unmarshaling: %s", err)
		}
		if tt.empty {
			if tc.UpdateStrategy.RegistryPoll.Interval != nil {
				t.Fatal("expected nil interval")
			}
			continue
		}
		if tc.UpdateStrategy.RegistryPoll.Interval.String() != tt.out.RegistryPoll.Interval.String() {
			t.Fatalf("expected %s, got %s", tt.out.RegistryPoll.Interval.String(), tc.UpdateStrategy.RegistryPoll.Interval.String())
		}
	}
}
