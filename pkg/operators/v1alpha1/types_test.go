package v1alpha1

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetAllCRDDescriptions(t *testing.T) {
	var table = []struct {
		owned    []string
		required []string
		expected []string
	}{
		{nil, nil, nil},
		{[]string{}, []string{}, []string{}},
		{[]string{"owned"}, []string{}, []string{"owned"}},
		{[]string{}, []string{"required"}, []string{"required"}},
		{[]string{"owned"}, []string{"required"}, []string{"owned", "required"}},
		{[]string{"first", "second"}, []string{"first", "second"}, []string{"first", "second"}},
		{[]string{"first", "second", "third"}, []string{"second", "third", "fourth"}, []string{"first", "second", "third", "fourth"}},
	}

	for _, tt := range table {
		// Build a list of owned CRDDescription used in the CSV.
		ownedDescriptions := make([]CRDDescription, 0)
		for _, crdName := range tt.owned {
			ownedDescriptions = append(ownedDescriptions, CRDDescription{
				Name: crdName,
			})
		}

		// Build a list of owned CRDDescription used in the CSV.
		requiredDescriptions := make([]CRDDescription, 0)
		for _, crdName := range tt.required {
			requiredDescriptions = append(requiredDescriptions, CRDDescription{
				Name: crdName,
			})
		}

		// Build a list of expected CRDDescriptions.
		expectedDescriptions := make([]CRDDescription, 0)
		sort.StringSlice(tt.expected).Sort()
		for _, expectedName := range tt.expected {
			expectedDescriptions = append(expectedDescriptions, CRDDescription{
				Name: expectedName,
			})
		}

		// Create a blank CSV with the owned descriptions.
		csv := ClusterServiceVersion{
			Spec: ClusterServiceVersionSpec{
				CustomResourceDefinitions: CustomResourceDefinitions{
					Owned:    ownedDescriptions,
					Required: requiredDescriptions,
				},
			},
		}

		// Call GetAllCRDDescriptions and ensure the result is as expected.
		require.Equal(t, expectedDescriptions, csv.GetAllCRDDescriptions())
	}
}

func TestOwnsCRD(t *testing.T) {
	var table = []struct {
		ownedCRDNames []string
		crdName       string
		expected      bool
	}{
		{nil, "", false},
		{nil, "querty", false},
		{[]string{}, "", false},
		{[]string{}, "querty", false},
		{[]string{"owned"}, "owned", true},
		{[]string{"owned"}, "notOwned", false},
		{[]string{"first", "second"}, "first", true},
		{[]string{"first", "second"}, "second", true},
		{[]string{"first", "second"}, "third", false},
	}

	for _, tt := range table {
		// Build a list of CRDDescription used in the CSV.
		var ownedDescriptions []CRDDescription
		for _, crdName := range tt.ownedCRDNames {
			ownedDescriptions = append(ownedDescriptions, CRDDescription{
				Name: crdName,
			})
		}

		// Create a blank CSV with the owned descriptions.
		csv := ClusterServiceVersion{
			Spec: ClusterServiceVersionSpec{
				CustomResourceDefinitions: CustomResourceDefinitions{
					Owned: ownedDescriptions,
				},
			},
		}

		// Call OwnsCRD and ensure the result is as expected.
		require.Equal(t, tt.expected, csv.OwnsCRD(tt.crdName))
	}
}

func TestCatalogSource_Update(t *testing.T) {
	var table = []struct {
		description string
		catsrc      CatalogSource
		result      bool
		sleep       time.Duration
	}{
		{
			description: "polling interval set: last update time zero: update for the first time",
			catsrc: CatalogSource{
				ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Time{Time: time.Now()}},
				Spec: CatalogSourceSpec{
					UpdateStrategy: &UpdateStrategy{
						RegistryPoll: &RegistryPoll{
							Interval: &metav1.Duration{Duration: 1 * time.Second},
						},
					},
					Image:      "mycatsrcimage",
					SourceType: SourceTypeGrpc},
			},
			result: true,
			sleep:  2 * time.Second,
		},
		{
			description: "polling interval set: time to update based on previous poll timestamp",
			catsrc: CatalogSource{
				ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Time{Time: time.Now()}},
				Spec: CatalogSourceSpec{
					UpdateStrategy: &UpdateStrategy{
						RegistryPoll: &RegistryPoll{
							Interval: &metav1.Duration{Duration: 1 * time.Second},
						},
					},
					Image:      "mycatsrcimage",
					SourceType: SourceTypeGrpc,
				},
				Status: CatalogSourceStatus{LatestImageRegistryPoll: &metav1.Time{Time: time.Now()}},
			},
			result: true,
			sleep:  2 * time.Second,
		},
		{
			description: "polling interval set: not time to update based on previous poll timestamp",
			catsrc: CatalogSource{
				ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Time{Time: time.Now()}},
				Spec: CatalogSourceSpec{
					UpdateStrategy: &UpdateStrategy{
						RegistryPoll: &RegistryPoll{
							Interval: &metav1.Duration{Duration: 1 * time.Second},
						},
					},
					Image:      "mycatsrcimage",
					SourceType: SourceTypeGrpc,
				},
				Status: CatalogSourceStatus{LatestImageRegistryPoll: &metav1.Time{Time: time.Now()}},
			},
			result: true,
			sleep:  2 * time.Second,
		},
	}

	for i, tt := range table {
		time.Sleep(table[i].sleep)
		require.Equal(t, tt.result, table[i].catsrc.Update(), table[i].description)
	}
}

func TestCatalogSource_Poll(t *testing.T) {
	var table = []struct {
		description string
		catsrc      CatalogSource
		result      bool
	}{
		{
			description: "poll interval set to zero: do not check for updates",
			catsrc:      CatalogSource{Spec: CatalogSourceSpec{}},
			result:      false,
		},
		{
			description: "not image based catalog source: do not check for updates",
			catsrc: CatalogSource{Spec: CatalogSourceSpec{SourceType: SourceTypeInternal,
				Address: "127.0.0.1:8080"}},
			result: false,
		},
		{
			description: "polling set with image based catalog: check for updates",
			catsrc: CatalogSource{Spec: CatalogSourceSpec{
				Image:      "my-image",
				SourceType: SourceTypeGrpc,
				UpdateStrategy: &UpdateStrategy{
					RegistryPoll: &RegistryPoll{
						Interval: &metav1.Duration{Duration: 1 * time.Second},
					},
				},
			},
			},
			result: true,
		},
	}
	for i, tt := range table {
		require.Equal(t, tt.result, table[i].catsrc.Poll(), table[i].description)
	}
}

func TestUpdateStrategyUnmarshal(t *testing.T) {
	type TestStruct struct {
		UpdateStrategy UpdateStrategy `json:"updateStrategy,omitempty"`
	}
	validDuration, err := time.ParseDuration("45m")
	if err != nil {
		panic(fmt.Errorf("error parsing duration: %s", err))
	}
	defaultDuration, err := time.ParseDuration("15m")
	if err != nil {
		panic(fmt.Errorf("error parsing duration: %s", err))
	}
	tests := []struct {
		name string
		in   []byte
		out  TestStruct
		err  error
	}{
		{
			name: "valid",
			in:   []byte(`{"UpdateStrategy": {"registryPoll":{"interval":"45m"}}}`),
			out: TestStruct{
				UpdateStrategy{
					&RegistryPoll{
						RawInterval:  "45m",
						Interval:     &metav1.Duration{Duration: validDuration},
						ParsingError: "",
					},
				},
			},
		},
		{
			name: "invalid",
			in:   []byte(`{"UpdateStrategy": {"registryPoll":{"interval":"19mError Code"}}}`),
			out: TestStruct{
				UpdateStrategy{
					&RegistryPoll{
						RawInterval:  "19mError Code",
						Interval:     &metav1.Duration{Duration: defaultDuration},
						ParsingError: "error parsing spec.updateStrategy.registryPoll.interval. Using the default value of 15m0s instead. Error: time: unknown unit \"mError Code\" in duration \"19mError Code\"",
					},
				},
			},
		},
		{
			name: "empty",
			in:   []byte(`{"UpdateStrategy": {"registryPoll":{"interval":""}}}`),
			out: TestStruct{
				UpdateStrategy{
					&RegistryPoll{
						Interval:     &metav1.Duration{Duration: defaultDuration},
						ParsingError: "error parsing spec.updateStrategy.registryPoll.interval. Using the default value of 15m0s instead. Error: time: invalid duration \"\"",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TestStruct{}
			err := json.Unmarshal(tt.in, &s)
			require.Equal(t, tt.out.UpdateStrategy.RawInterval, s.UpdateStrategy.RawInterval)
			require.Equal(t, tt.out.UpdateStrategy.Interval, s.UpdateStrategy.Interval)
			require.Equal(t, tt.out.UpdateStrategy.ParsingError, s.UpdateStrategy.ParsingError)
			require.Equal(t, tt.err, err)
		})
	}
}
